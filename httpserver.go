package main

import (
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

func runHTTPServer(host string, port int, workerId uint64, datacenterId uint64) {
	worker, err := idworker.NewIdWorker(workerId, datacenterId)

	if err == nil {
		http.HandleFunc("/api/snowflake", func(w http.ResponseWriter, r *http.Request) {
			snowflakeHandler(w, r, worker)
		})

		http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	} else {
		log.Fatal(err)
	}
}
