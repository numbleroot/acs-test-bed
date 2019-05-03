.PHONY: all clean genconfigs runexperiments build syncscripts

all: clean build

clean:
	go clean -i ./...
	rm -rf genconfigs runexperiments

genconfigs:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/genconfigs

runexperiments:
	CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' ./cmd/runexperiments

build: genconfigs runexperiments

syncscripts:
	gsutil cp scripts/startup.sh gs://acs-eval/
	gsutil cp scripts/zeno-pki_eval.sh gs://acs-eval/
	gsutil cp scripts/zeno_client_eval.sh gs://acs-eval/
	gsutil cp scripts/zeno_mix_eval.sh gs://acs-eval/
	gsutil cp scripts/vuvuzela-client_eval.sh gs://acs-eval/
	gsutil cp scripts/vuvuzela-coordinator_eval.sh gs://acs-eval/
	gsutil cp scripts/vuvuzela-mixer_eval.sh gs://acs-eval/
	gsutil cp scripts/pung_client_eval.sh gs://acs-eval/
	gsutil cp scripts/pung_server_eval.sh gs://acs-eval/

cleanresults:
	rm -rf results/*
