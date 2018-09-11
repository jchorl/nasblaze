build:
	docker build -t jchorl/nasblaze .

run:
	docker run -it --rm jchorl/nasblaze

pi:
	GOOS=linux GOARCH=arm GOARM=5 go build ./...

default: build run

.PHONY: build run
