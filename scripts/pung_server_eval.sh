#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -server -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Configure tc according to environment variable.

# Add iptables rules to be able to count number of transferred
# bytes for evaluation and initialize them to zero.
iptables -A INPUT -p tcp --dport 33000
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run pung as server.
/root/pung-server -e 35 -n 1 -w 1 -i ${LISTEN_IP} -s 33000 -k 1 -t t -d 2 -b 0 -m 1000 > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/
