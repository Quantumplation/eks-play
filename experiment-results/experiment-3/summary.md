Containers: 1
Server type: Go
Total requests: 3 million
Parallel Requests: 500 per container
EOF Errors: 0
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 0
ECONNREFUSED: 0

The goal of this experiment was to scale down the previous one.  With one container, it's only talking to itself.  Even with high parallelism, we saw no ECONNRESETs.

Got concerned that it might be a duration thing, so I re-ran this experiment again for over 24 hours.

Containers: 1
Server type: Go
Total requests: 266 million
Parallel Requests: 500 per container
EOF Errors: 0
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 0
ECONNREFUSED: 0