.PHONY: all clean genconf runexp build syncscripts

all: clean build

clean:
	go clean -i ./...
	rm -rf genconfigs runexperiments

genconf:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

runexp:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

build: genconf runexp

syncscripts:
	gsutil cp scripts/startup.sh gs://acs-eval/
	gsutil cp scripts/zeno-pki_eval.sh gs://acs-eval/
	gsutil cp scripts/zeno_client_eval.sh gs://acs-eval/