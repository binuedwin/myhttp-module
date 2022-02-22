all: test build

build:
	go build -o myhttp main.go

test:
	go test -v