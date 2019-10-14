package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Statistics ...
type Statistics struct {
	hostname string

	totalOutgoingRequests      int
	successfulOutgoingRequests int
	failedOutgoingRequests     int
	outgoingHTTPErrors         int
	outgoingNetworkErrors      int
	outgoingUnknownErrors      int
	econnresetErrors           int

	totalIncomingRequests int
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Printf("%s started...", hostname)

	host := os.Getenv("GO_SERVICE_SERVICE_HOST")
	port := os.Getenv("GO_SERVICE_SERVICE_PORT")

	stats := Statistics{}
	stats.hostname = hostname

	baseURL := fmt.Sprintf("http://%s:%s", host, port)
	if host == "" || port == "" {
		baseURL = "http://localhost:8080"
	}
	url := fmt.Sprintf("%s/sample", baseURL)

	log.Printf("Continually requesting: %s", url)

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		stats.totalIncomingRequests++
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
		stats.totalOutgoingRequests++
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			stats.failedOutgoingRequests++
			stats.outgoingNetworkErrors++
			continue
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			stats.failedOutgoingRequests++
			stats.outgoingUnknownErrors++
			continue
		}
		if resp.StatusCode != 200 {
			stats.failedOutgoingRequests++
			stats.outgoingHTTPErrors++
		}
		stats.successfulOutgoingRequests++
		resp.Body.Close()
	}
}
