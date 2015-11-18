default: build test

build:
	CGO_ENABLED=0 go build -ldflags "-s" -a -installsuffix cgo -o gosnowflake .

test:
	go test
