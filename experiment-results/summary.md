The investigation into ECONNRESET was composed of a series of experiments.

First, I spun up an EKS cluster running kubernetes 1.14.  I scaled the nodes up and down for each experiment, usually targeting the ability to run between 600 and 700 containers.

The first batch of experiments consisted of containers running a simple go service I wrote.  This service:
 - Started an HTTP service listening for requests to /sample, to which it responded with a 200
 - Used environment variable injected routes to hit other instances of the service
 - Recorded request rate and error rate statistics
 - Saved statistics and unrecognized errors to dynamodb

The logic for tracking the request statistics (found in the statistics.csv files) is as follows:
 - Each incoming request incremented `TotalIncomingRequests`
 - Each outgoing request incremented `TotalOutgoingRequests`
 - If we encountered an error (detected in any of 3 ways), this incremented `FailedOutgoingRequests`
   - If the go http client returned an err, this incremented `OutgoingNetworkErrors`, and then tried to classify the error
     - If the err was a wrapper for `io.EOF`, this incremented `EOFErrors`.
     - If the err was a wrapper for `syscall.ECONNRESET`, this incremented `TrueECONNRESETErrors`
     - If the err was a wrapper for `syscall.ECONNREFUSED`, this incremented `ECONNREFUSEDErrors`
     - If the err was a wrapper for `syscall.ECONNABORTED`, this incremented `ECONNABORTEDErrors`
     - Otherwise, the error was logged in a different dynamodb table, found in errors.csv
   - If we encountered an error while reading the response body, this incremented `OutgoingUnknownErrors`
   - If we received a non-200 HTTP status code, this incremented `OutgoingHTTPErrors`
 - Otherwise, `SuccessfulOutgoingRequests` was incremented.

 I'd like to highlight the specific colleciton of errors I was manually able to reproduce.  These can be found in the `Errors` folder, and each consists of some kind of server, and `socket-client.go` designed to connect to any of them and reproduce the error:
  - `syscall.ECONNREFUSED`: These errors occured when there was nothing listening on the accepting socket (`Errors/wrong-port.go`)
  - `syscall.ECONNRESET`: These errors occured when the server accepts the connection, and immediately closes it `Errors/graceful-close.go`
  - `io.EOF`: These errors occured when I wrote a server that accepted a socket connection and immediately exited the process (`Errors/sudden-terminate.go`)
  - `syscall.ECONNABORTED`: I found this while reading golang syscall error codes, but was unable to reproduce it

We also received these errors in our experiment that I was unable to figure out how to reproduce
 - `i/o timeout`: I tried to just leave a socket hanging (`Errors/hanging.go`), but this never seemed to time out
 - `no route to host`: I imagine this has something to do with DNS resolution or proxy states, but I don't have a way to reproduce this.


 # Experiment 1

  - Low parallelism, but massive scale.  0 ECONNRESET errors of any kind, despite doing half a billion network requests

# Experiment 2

  - Increased the parallelism, so each service could make up to 500 requests in parallel.  SAW ECONNRESET ERRORS!

# Experiment 3

  - Ran a single container to try to reproduce the ECONNRESETs