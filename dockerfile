FROM golang:alpine
EXPOSE 8080


RUN go build *.go
CMD ["./command-service"]
