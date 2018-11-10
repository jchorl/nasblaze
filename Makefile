build:
	docker build -t jchorl/nasblaze .

run:
	docker run -it --rm jchorl/nasblaze

pi:
	GOOS=linux GOARCH=arm GOARM=5 go build -o nasblaze main.go

default: build run

.PHONY: build run
