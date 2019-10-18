.PHONY: all clean calcstats collector genconfigs operator runexperiments build syncbucket

all: clean build

clean:
	go clean -i ./...
	rm -rf calcstats collector genconfigs operator runexperiments

build: calcstats collector genconfigs operator runexperiments

calcstats:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/calcstats

collector:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/collector

genconfigs:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

operator:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/operator

runexperiments:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments
