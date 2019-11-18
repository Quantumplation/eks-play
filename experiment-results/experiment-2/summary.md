Containers: 381 (~500, but some failed to write their statistics to dynamodb)
Server type: Go
Total requests: 55 million
Parallel Requests: 500 per container
EOF Errors: 4276
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 1627
ECONNREFUSED: 4211

The goal of this experiment was to determine if parallelism had anything to do with the ECONNRESETs, which seems to be the case!