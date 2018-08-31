FROM golang:1.11
WORKDIR /go/src/github.com/jchorl/nasblaze
COPY main.go .
RUN go get ./... && GOOS=linux go build -o nasblaze .
FROM golang:1.11

FROM ubuntu:18.04
RUN apt-get update && \
    apt-get install -y curl unzip man-db && \
    curl https://rclone.org/install.sh | bash
WORKDIR /nasblaze
COPY rclone.conf .
COPY --from=0 /go/src/github.com/jchorl/nasblaze/nasblaze .
CMD ["./nasblaze", "-logtostderr=true"]
