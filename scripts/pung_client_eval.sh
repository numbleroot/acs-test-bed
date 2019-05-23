#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -client -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Configure tc according to environment variable.

# Add iptables rules to be able to count number of transferred
# bytes for evaluation and initialize them to zero.
iptables -A INPUT -m state --state ESTABLISHED,RELATED
iptables -A OUTPUT -p tcp --dport 33000
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run pung as client.
/root/pung-client -e /tmp/collect -n "${NAME_OF_NODE}" -p "${PARTNER_OF_NODE}" \
    -x "ACS_SECRET_${NAME_OF_NODE}_${PARTNER_OF_NODE}" -h ${PKI_IP}:33000 \
    -r 35 -k 1 -s 1 -t t -d 2 -b 0 > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/
