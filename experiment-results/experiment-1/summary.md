Containers: 1267 (~600, with a rolling restart in the middle)
Server type: Go
Total requests: 547 million
Parallel Requests: 1 per container
EOF Errors: 0
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 0
ECONNREFUSED: 2437

The goal of this was to mainly get all the infrastructure set up, and prove that it is possible to live in a world without ECONNRESETs.