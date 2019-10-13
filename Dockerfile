FROM alpine:latest

COPY go-server.exe .

CMD ["go-server.exe"]