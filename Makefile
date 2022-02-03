GOOS=linux
GOARCH=arm
GOARM=7
export GOOS GOARCH GOARM

all: generate build

generate:
	@echo "Generating..."
	go build github.com/porjo/youtubeuploader

build:
	@echo "Building..."
	go build .
