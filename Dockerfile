FROM golang:1.8 AS build
WORKDIR /
COPY go-server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server go-server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=build /server .
ENTRYPOINT ["./server"]
