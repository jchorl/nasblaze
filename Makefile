build:
	docker build -t jchorl/nasblaze .

run:
	docker run -it --rm jchorl/nasblaze

default: build run
