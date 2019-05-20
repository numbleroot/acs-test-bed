#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -client -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Configure tc according to environment variable.

# Reset bytes counters right before starting pung.
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run pung as client.
/root/pung -n 0 -p 0 -x "Shared_Secret_ACS_Eval" -h ${PKI_IP}:33000 -d 2 -r 25 -b 0 -k 64 -o h2 -t b > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/
