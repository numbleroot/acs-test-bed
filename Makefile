.PHONY: all clean genconfigs collector runexperiments build syncbucket

all: clean build

clean:
	go clean -i ./...
	rm -rf genconfigs collector runexperiments

build: genconfigs collector runexperiments

genconfigs:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

collector:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/collector

runexperiments:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

syncbucket:
	gsutil -m cp scripts/* gs://acs-eval/
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno/zeno gs://acs-eval/zeno
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno-pki/zeno-pki gs://acs-eval/zeno-pki
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno-eval/collector gs://acs-eval/collector

cleanresults:
	rm -rf results/*
