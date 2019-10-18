package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

// ENABLEDYNAMO Whether to enable dynamo
const ENABLEDYNAMO = true

// Statistics ...
type Statistics struct {
	Hostname string

	TotalOutgoingRequests      int32
	SuccessfulOutgoingRequests int32
	FailedOutgoingRequests     int32
	OutgoingHTTPErrors         int32
	OutgoingNetworkErrors      int32
	OutgoingUnknownErrors      int32
	EOFErrors                  int32
	TrueECONNRESETErrors       int32
	ECONNREFUSEDErrors         int32
	ECONNABORTEDErrors         int32
	ForceClosedErrors          int32

	TotalIncomingRequests int32
}

func updateStats(stats *Statistics, lock *sync.RWMutex) {
	counter := 0
	for {
		lock.Lock()
		b, _ := json.MarshalIndent(stats, "  ", "\t")
		av, _ := dynamodbattribute.MarshalMap(stats)
		lock.Unlock()

		log.Print("Statistics:")
		log.Printf("%s", string(b))
		if ENABLEDYNAMO && counter%20 == 0 {
			log.Print("Saving...")
			sess := session.Must(session.NewSession(&aws.Config{
				Region: aws.String("us-east-1")},
			))
			svc := dynamodb.New(sess)
			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String("eks-play-statistics"),
			}
			_, err := svc.PutItem(input)
			if err != nil {
				log.Printf("Couldn't update statistics: %v", err)
			} else {
				log.Print("Saved")
			}
		}
		counter++
		time.Sleep(5 * time.Second)
	}
}

func recordError(host string, e error) {
	// Go through a lot of effort to preserve as much info as we can
	eb, _ := json.Marshal(e)
	log.Printf("Unrecognized error: %s", string(eb))
	log.Print(fmt.Errorf("Unrecognized error: %w", e))
	if ENABLEDYNAMO {

		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		))
		svc := dynamodb.New(sess)

		av, _ := dynamodbattribute.MarshalMap(e)
		id, _ := uuid.NewRandom()

		h, _ := dynamodbattribute.Marshal(host)
		idm, _ := dynamodbattribute.Marshal(id.String())
		ej, _ := dynamodbattribute.Marshal(string(eb))
		et, _ := dynamodbattribute.Marshal(e.Error())

		av["Id"] = idm
		av["Host"] = h
		av["ErrorJson"] = ej
		av["ErrorText"] = et

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("eks-play-errors"),
		}
		_, err := svc.PutItem(input)
		if err != nil {
			log.Printf("Couldn't insert error: %v", err)
		}
	}
}

type result struct {
	resp *http.Response
	err  error
}

func doResponseLoop(ch chan result, stats *Statistics, lock *sync.RWMutex) {
	for {
		res := <-ch
		resp := res.resp
		err := res.err
		lock.RLock()
		if err != nil {
			atomic.AddInt32(&stats.FailedOutgoingRequests, 1)
			atomic.AddInt32(&stats.OutgoingNetworkErrors, 1)
			if errors.Is(err, io.EOF) {
				atomic.AddInt32(&stats.EOFErrors, 1)
			} else if errors.Is(err, syscall.ECONNRESET) {
				atomic.AddInt32(&stats.TrueECONNRESETErrors, 1)
			} else if errors.Is(err, syscall.ECONNREFUSED) {
				atomic.AddInt32(&stats.ECONNREFUSEDErrors, 1)
			} else if errors.Is(err, syscall.ECONNABORTED) {
				atomic.AddInt32(&stats.ECONNABORTEDErrors, 1)
			} else {
				recordError(stats.Hostname, err)
			}
			lock.RUnlock()
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			atomic.AddInt32(&stats.FailedOutgoingRequests, 1)
			atomic.AddInt32(&stats.OutgoingUnknownErrors, 1)
			lock.RUnlock()
			continue
		}
		if resp.StatusCode != 200 {
			log.Print(string(b))
			atomic.AddInt32(&stats.FailedOutgoingRequests, 1)
			atomic.AddInt32(&stats.OutgoingHTTPErrors, 1)
			lock.RUnlock()
			continue
		}
		atomic.AddInt32(&stats.SuccessfulOutgoingRequests, 1)
		lock.RUnlock()
		resp.Body.Close()
		runtime.Gosched()
	}
}

func doRequestLoop(url string, ch chan result, stats *Statistics, lock *sync.RWMutex) {
	for {
		lock.RLock()
		atomic.AddInt32(&stats.TotalOutgoingRequests, 1)
		lock.RUnlock()
		resp, err := http.Get(url)
		ch <- result{resp, err}
		runtime.Gosched()
	}
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Printf("%s started...", hostname)

	lock := sync.RWMutex{}
	stats := Statistics{}
	stats.Hostname = hostname

	host := os.Getenv("GO_SERVICE_SERVICE_HOST")
	port := os.Getenv("GO_SERVICE_SERVICE_PORT")

	baseURL := fmt.Sprintf("http://%s:%s", host, port)
	if host == "" || port == "" {
		baseURL = "http://localhost:8080"
	}
	url := fmt.Sprintf("%s/sample", baseURL)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		lock.RLock()
		atomic.AddInt32(&stats.TotalIncomingRequests, 1)
		lock.RUnlock()
		fmt.Fprintf(w, "")
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		b, _ := json.Marshal(stats)
		lock.Unlock()
		w.Write(b)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Print("Sleeping for some random time to let things warm up...")
	time.Sleep(20*time.Second + time.Duration(rand.Intn(60))*time.Second)
	log.Printf("Continually requesting: %s", url)

	go updateStats(&stats, &lock)

	parallelism := 300
	ch := make(chan result, parallelism)

	for i := 0; i < parallelism; i++ {
		go doRequestLoop(url, ch, &stats, &lock)
	}
	go doResponseLoop(ch, &stats, &lock)

	for {
		time.Sleep(1 * time.Minute)
	}
}
