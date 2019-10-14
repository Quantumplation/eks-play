package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Statistics ...
type Statistics struct {
	Hostname string

	TotalOutgoingRequests      int32
	SuccessfulOutgoingRequests int32
	FailedOutgoingRequests     int32
	OutgoingHTTPErrors         int32
	OutgoingNetworkErrors      int32
	OutgoingUnknownErrors      int32
	EconnresetErrors           int32

	TotalIncomingRequests int32
}

func updateStats(stats *Statistics, lock *sync.RWMutex) {
	for {
		time.Sleep(1 * time.Minute)
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		))
		svc := dynamodb.New(sess)
		lock.Lock()
		av, _ := dynamodbattribute.MarshalMap(stats)
		lock.Unlock()
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("eks-play-statistics"),
		}
		_, err := svc.PutItem(input)
		if err != nil {
			log.Printf("Couldn't update statistics: %v", err)
		}
	}
}

func printStats(stats *Statistics, lock *sync.RWMutex) {
	for {
		lock.Lock()
		b, _ := json.MarshalIndent(stats, "  ", "\t")
		lock.Unlock()
		log.Print("Statistics: ")
		log.Printf("%s", string(b))
		time.Sleep(5 * time.Second)
	}
}

func doRequestLoop(url string, stats *Statistics, lock *sync.RWMutex) {
	for {
		lock.RLock()
		atomic.AddInt32(&stats.TotalOutgoingRequests, 1)
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			atomic.AddInt32(&stats.FailedOutgoingRequests, 1)
			atomic.AddInt32(&stats.OutgoingNetworkErrors, 1)
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
		time.Sleep(0)
	}
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Printf("%s started...", hostname)

	host := os.Getenv("GO_SERVICE_SERVICE_HOST")
	port := os.Getenv("GO_SERVICE_SERVICE_PORT")

	lock := sync.RWMutex{}
	stats := Statistics{}
	stats.Hostname = hostname
	go printStats(&stats, &lock)
	go updateStats(&stats, &lock)

	baseURL := fmt.Sprintf("http://%s:%s", host, port)
	if host == "" || port == "" {
		baseURL = "http://localhost:8080"
	}
	url := fmt.Sprintf("%s/sample", baseURL)

	log.Printf("Continually requesting: %s", url)

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&stats.TotalIncomingRequests, 1)
		fmt.Fprintf(w, "Good")
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

	for i := 0; i < 1; i++ {
		go doRequestLoop(url, &stats, &lock)
	}
	for {
		time.Sleep(1 * time.Minute)
	}
}
