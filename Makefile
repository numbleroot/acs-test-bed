.PHONY: all clean build

all: clean build

clean:
	go clean -i ./...
	rm -rf genconfigs runexperiments

genconf:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

runexp:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

build: genconf runexp