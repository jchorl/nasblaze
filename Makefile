build:
	docker build -t jchorl/nasblaze .

run:
	docker run -it --rm jchorl/nasblaze

pi:
	docker container run --rm -it \
		-v $(PWD):/nasblaze \
		-w /nasblaze \
		-e GOOS=linux \
		-e GOARCH=arm \
		-e GOARM=5 \
		golang:1.11.2 \
		go build -o nasblaze main.go

default: build run

.PHONY: build run
