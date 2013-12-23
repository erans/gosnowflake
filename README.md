gosnowflake
=========

Twitter's Snowflake server implemented in Go.

Currently supports a simple HTTP interface. Thrift interface conforming to Twitter's snowflake will follow soon.

Execute this curl command to get a new ID:
```
curl http://localhost:8080/api/snowflake
```