package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/erans/gosnowflake/idworker"
	"github.com/erans/gosnowflake/snowflake"
	"log"
	"os"
)

func runThriftServer(host string, port int, protocol string, framed bool, buffered bool, workerId uint64, datacenterId uint64) {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	var err error
	var transport thrift.TServerTransport
	transport, err = thrift.NewTServerSocket(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	worker, err := idworker.NewIdWorker(workerId, datacenterId)

	processor := snowflake.NewSnowflakeProcessor(worker)

	if err == nil {
		server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
		server.Serve()
	} else {
		log.Fatal(err)
	}
}
