Containers: 10
Server type: Go
Total requests: 4 million
Parallel Requests: 500 per container
EOF Errors: 7
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 61
ECONNREFUSED: 267

The goal of this experiment was to confirm ECONNRESETs again before trying advice from EKS SME.  Wanted to make sure I could still spin up the cluster and reproduce etc.