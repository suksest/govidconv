BINARY=engine
all: setup test clean run

setup: 
	sudo apt-get install ffmpeg

test : 
		go test ./...

build:
		go build -o ${BINARY} app/main.go

clean:
		@echo "cleaning built apps from local storage..."
		@if [ -f ${BINARY} ] ; then rm -f ${BINARY} ; fi

docker:
		docker build -t govidconv:latest .

run:
		docker-compose up --build -d

stop:
		docker-compose down

.PHONY: test clean build docker run stop