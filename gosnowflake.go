package main

import (
	"flag"
	"fmt"
	"github.com/erans/gosnowflake/idworker"
	"log"
	"net/http"
)

func snowflakeHandler(w http.ResponseWriter, r *http.Request, worker *idworker.IdWorker) {
	id, err := worker.Next()
	if err == nil {
		fmt.Fprintf(w, "%d", id)
	}
}

var (
	host         string
	port         int
	workerId     int
	datacenterId int
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "The host to listen to (default: 0.0.0.0")
	flag.IntVar(&port, "port", 8080, "The port to listen on (default for http: 8080, thrift: 7609)")
	flag.IntVar(&workerId, "workerid", 1, "The worker ID. Use different ones for different processes (default: 1)")
	flag.IntVar(&datacenterId, "datacenterid", 1, "The data center ID. Use different ones for different data center (default: 1)")
}

func main() {
	flag.Parse()

	worker, err := idworker.NewIdWorker(1, 1)

	if err == nil {
		http.HandleFunc("/api/snowflake", func(w http.ResponseWriter, r *http.Request) {
			snowflakeHandler(w, r, worker)
		})

		http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	} else {
		log.Fatal(err)
	}
}
