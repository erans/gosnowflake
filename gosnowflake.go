package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	serverType     string
	host           string
	port           int
	thriftProtocol string
	thriftFramed   bool
	thriftBuffered bool
	workerId       int
	datacenterId   int
)

func Usage() {
	fmt.Fprint(os.Stderr, "\ngosnowflake - a Go based implementation of Twitter's Snowflake unique ID generation service\n")
	fmt.Fprint(os.Stderr, "Written by Eran Sandler (@erans)\n")
	fmt.Fprint(os.Stderr, "Follow and contribute at: https://github.com/erans/gosnowflake\n\n")
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func main() {
	flag.Usage = Usage
	serverType := flag.String("servertype", "http", "The type of server to run: http or thrift (default: http)")
	host := flag.String("host", "0.0.0.0", "The host to listen to (default: 0.0.0.0)")
	port := flag.Int("port", 8080, "The port to listen on (default for http: 8080, thrift: 7609)")
	thriftProtocol := flag.String("thriftprotocol", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	thriftFramed := flag.Bool("thriftframed", false, "Use framed transport")
	thriftBuffered := flag.Bool("thriftbuffered", false, "Use buffered transport")
	workerId := flag.Int("workerid", 1, "The worker ID. Use different ones for different processes (default: 1)")
	datacenterId := flag.Int("datacenterid", 1, "The data center ID. Use different ones for different data center (default: 1)")

	flag.Parse()

	switch *serverType {
	case "http":
		runHTTPServer(*host, *port, uint64(*workerId), uint64(*datacenterId))
	case "thrift":
		runThriftServer(*host, *port, *thriftProtocol, *thriftFramed, *thriftBuffered, uint64(*workerId), uint64(*datacenterId))
	default:
		fmt.Fprint(os.Stderr, "Invalid servertype specified", serverType, "\n")
		Usage()
		os.Exit(1)
	}
}
