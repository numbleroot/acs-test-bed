#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -client -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Determine active network device.
NET_DEVICE=$(ip addr | awk '/state UP/ {print $2}' | sed 's/.$//')
printf "Found active network device: '${NET_DEVICE}'.\n"

# Configure tc according to environment variable.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc add dev ${NET_DEVICE} root ${TC_CONFIG}
    printf "Configured ${NET_DEVICE} with tc parameters.\n"
fi

# Add iptables rules to be able to count number of transferred
# bytes for evaluation and initialize them to zero.
iptables -A INPUT -m state --state ESTABLISHED,RELATED
iptables -A OUTPUT -p tcp --dport 33000
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run pung as client.
/root/pung-client -e /tmp/collect -n "${NAME_OF_NODE}" -p "${PARTNER_OF_NODE}" \
    -x "ACS_EVALUATION_SECRET" -h ${PKI_IP}:33000 -r 30 -k 1 -s 1 -t b -d 2 -b 0 > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc del dev ${NET_DEVICE} root
fi

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/
