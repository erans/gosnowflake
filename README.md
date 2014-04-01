# gosnowflake

Twitter's Snowflake server implemented in Go.
Written by Eran Sandler ([@erans](http://twitter.com/erans))

Supports two interfaces:
- Simple HTTP -
Execute this curl command to get a new ID:
```
curl http://localhost:8080/api/snowflake
```

- Thrift - Initial support using Thrift 0.9.1.
    - Supports binary, compact, json and simple json protocols
    - Framed and buffered transport
    - No support for TLS

## Build
Just grab it and run ```go build```. I've built and tested it with Go 1.2 on Linux and Mac.

## Benefits over other implementations
- No dependency on any runtime, just deploy the executable and run it.
- Performace should be near native as Go is a fully compiled language

## Performance

### HTTP Interface

#### Methodology

##### Hardware
Amazon Web Services (AWS) c3.large (their newest to date high cpu instance type) with 2 vCPU

##### Test
Tests were performed using Apache Bench (`ab`) tool with the following parameters:
```
ab -n 100000 -c X http://server/api/snowflake
```
So we are 100,000 requets and used concurrency X (see results table below for levels of concurrency checked).

`ab` was run from a different t1.micro instance residing in a different availability zone (AZ) than the test server


#### Results
| Concurrency | # of requests | Total time (sec) | Reqs/Sec | Avg. CPU | 90% percentile request time (ms) |
|:-----------:|:-------------:|-----------------:|---------:|---------:|---------------------------------:|
|1            |100,000        |313.665           | 318      | 5%       | 3                                |
|5            |100,000        |63.832            | 1566     | 20%      | 4                                |
|10           |100,000        |32.774            | 3051     | 35%      | 4                                |
|15           |100,000        |22.998            | 4348     | 48%      | 4                                |
|20           |100,000        |18.394            | 5436     | 56%      | 4                                |
|25           |100,000        |15.765            | 6343     | 60%      | 5                                |
|30           |100,000        |14.265            | 7009     | 70%      | 5                                |

Twitter's Snowflake ([performance requirements](https://github.com/twitter/snowflake#requirements)) state:
- minimum 10k ids per second per process
- response rate 2ms (plus network latency)

As you can see we can get up to 7k ids per second, however with more than double the latency. With lower concurrency (20) we get about half the ids (5k) in twice the latency (4ms). I can associate some of the latecy for the current implementation with the fact that its a full HTTP text request and not binary like Thrift's default transport protocol and to the various network latencies which may or may not occur within AWS (though numerous runs shows its more of an implementation than AWS network latency)

In any case I never saw if the official implementation from Twitter even achieved these performance requirements and was too lazy to make the official implementation build and work with whatever sbt and Scala version it should (it doesn't really get built out-of-the-box). Hmpf.

## TODO
- ~~Implement Thrift support~~ (initial support is committed)
- Consider implementing a simple binary TCP interface since we don't really need the whole Thrift interface, just a simple command to get an ID and and 2 respones, one OK with the ID and the other an error.
- The current implementation uses the Twitter Epoch which is set way into the future causing the number to be very large - should it be changed to something else? Does it matter much?
- Consider a different approch in which we use the worker ID internally to associate a worker (with a dedicated id per code) and use datacenterid to differenciate machines. That way we can remove the lock, or keep it and avoid contention. Keep in mind datacenterid is only 5 bits (32) so a smart allocation of workerid with datacenterid if more than 32 machines is required.

