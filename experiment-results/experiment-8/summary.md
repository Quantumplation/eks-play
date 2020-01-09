Containers: 10
Server type: Go
Total requests: 
Parallel Requests: 500 per container
EOF Errors: 
Force Closed: 
ECONNABORTED: 
syscall ECONNRESET: 
ECONNREFUSED: 

This experiment recreated the results on a 1.12 kubernetes cluster off of EKS.
The next one will try k8s 1.16, and since EKS doesn't support that version yet, I wanted a baseline to compare against.