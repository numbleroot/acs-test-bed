#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -mix -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Configure tc according to environment variable.

# Reset bytes counters right before starting pung.
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run pung as server.
/root/pung-server -e 35 -n 1 -w 1 -i ${LISTEN_IP} -s 33000 -k 1 -t t -d 2 -b 0 -m 10 > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/
