package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Statistics ...
type Statistics struct {
	Hostname string

	TotalOutgoingRequests      int
	SuccessfulOutgoingRequests int
	FailedOutgoingRequests     int
	OutgoingHTTPErrors         int
	OutgoingNetworkErrors      int
	OutgoingUnknownErrors      int
	EconnresetErrors           int

	TotalIncomingRequests int
}

func printStats(stats *Statistics) {
	for {
		b, _ := json.MarshalIndent(stats, "  ", "\t")
		log.Print("Statistics: ")
		log.Printf("%s", string(b))
		time.Sleep(5 * time.Second)
	}
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Printf("%s started...", hostname)

	host := os.Getenv("GO_SERVICE_SERVICE_HOST")
	port := os.Getenv("GO_SERVICE_SERVICE_PORT")

	stats := Statistics{}
	stats.Hostname = hostname
	go printStats(&stats)

	baseURL := fmt.Sprintf("http://%s:%s", host, port)
	if host == "" || port == "" {
		baseURL = "http://localhost:8080"
	}
	url := fmt.Sprintf("%s/sample", baseURL)

	log.Printf("Continually requesting: %s", url)

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		stats.TotalIncomingRequests++
		fmt.Fprintf(w, "Good")
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for {
		stats.TotalOutgoingRequests++
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			stats.FailedOutgoingRequests++
			stats.OutgoingNetworkErrors++
			continue
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			stats.FailedOutgoingRequests++
			stats.OutgoingUnknownErrors++
			continue
		}
		if resp.StatusCode != 200 {
			stats.FailedOutgoingRequests++
			stats.OutgoingHTTPErrors++
		}
		stats.SuccessfulOutgoingRequests++
		resp.Body.Close()
	}
}
