package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

type Statistics struct {
	hostname string

	total_outgoing_requests int
	successful_outgoing_requests int
	failed_outgoing_requests int
	outgoing_http_errors int
	outgoing_network_errors int
	outgoing_unknown_errors int
	econnreset_errors int

	total_incoming_requests int
}

func main() {
	hostname := os.Getenv("HOSTNAME")
	log.Printf("%s started...", hostname)

	host := os.Getenv("GO_SERVICE_SERVICE_HOST")
	port := os.Getenv("GO_SERVICE_SERVICE_PORT")

	stats := Statistics{}
	stats.hostname = hostname

	base_url := fmt.Sprintf("http://%s:%s", host, port)
	if host == "" || port == "" {
		base_url = "http://localhost:8080"
	}
	url := fmt.Sprintf("%s/sample", base_url)

	log.Printf("Continually requesting: %s", url)

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		stats.total_incoming_requests += 1
		fmt.Fprintf(w, "Good")
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(stats)
		w.Write(b)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for {
		stats.total_outgoing_requests += 1
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			stats.failed_outgoing_requests += 1
			stats.outgoing_network_errors += 1
			continue
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			stats.failed_outgoing_requests += 1
			stats.outgoing_unknown_errors += 1
			continue
		}
		if resp.StatusCode != 200 {
			stats.failed_outging_requests += 1
			stats.outgoing_http_errors += 1
		}
		stats.successful_outgoing_requests += 1
		resp.Body.Close()
	}
}
