.PHONY: all clean calcstats collector genconfigs runexperiments build syncbucket

all: clean build

clean:
	go clean -i ./...
	rm -rf calcstats collector genconfigs runexperiments

build: calcstats collector genconfigs runexperiments

calcstats:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/calcstats

collector:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/collector

genconfigs:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

runexperiments:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

syncbucket:
	gsutil -m cp scripts/* gs://acs-eval/
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno/zeno gs://acs-eval/zeno
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno-pki/zeno-pki gs://acs-eval/zeno-pki
	gsutil cp ~/Rust/pung/target/release/client gs://acs-eval/pung-client
	gsutil cp ~/Rust/pung/target/release/server gs://acs-eval/pung-server
	gsutil cp ${GOPATH}/src/github.com/numbleroot/zeno-eval/collector gs://acs-eval/collector

cleanplots:
	rm -rf results/*/clients/*.png
	rm -rf results/*/servers/*.png
