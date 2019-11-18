FROM golang:latest AS build
WORKDIR /
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/google/uuid
COPY go-server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server go-server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=build /server .
ENTRYPOINT ["./server"]
