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

syncbucket:
	gsutil -m cp scripts/* gs://acs-eval/
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno/zeno gs://acs-eval/zeno
	gsutil cp ${GOPATH}/src/github.com/numbleroot/vuvuzela/client gs://acs-eval/vuvuzela-client
	gsutil cp ${GOPATH}/src/github.com/numbleroot/vuvuzela/coordinator gs://acs-eval/vuvuzela-coordinator
	gsutil cp ${GOPATH}/src/github.com/numbleroot/vuvuzela/mix gs://acs-eval/vuvuzela-mix
	gsutil cp ~/Rust/pung/target/release/client gs://acs-eval/pung-client
	gsutil cp ~/Rust/pung/target/release/server gs://acs-eval/pung-server
	gsutil cp ${GOPATH}/src/github.com/numbleroot/acs-test-bed/operator gs://acs-eval/operator
	gsutil cp ${GOPATH}/src/github.com/numbleroot/acs-test-bed/collector gs://acs-eval/collector
