.PHONY: all clean genconfigs runexperiments build syncbucket

all: clean build

clean:
	go clean -i ./...
	rm -rf genconfigs runexperiments

genconfigs:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

runexperiments:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

build: genconfigs runexperiments

syncbucket:
	gsutil cp scripts/* gs://acs-eval/
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno/zeno gs://acs-eval/zeno
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno-pki/zeno-pki gs://acs-eval/zeno-pki

cleanresults:
	rm -rf results/*
