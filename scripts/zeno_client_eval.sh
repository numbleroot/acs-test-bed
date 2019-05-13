#!/usr/bin/env bash

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki.pem /root/cert.pem
chmod 0644 /root/cert.pem

# Configure tc according to environment variable.

# Run metrics collector sidecar in background.
/root/collector -client -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Reset bytes counters right before starting zeno.
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run zeno as client.
/root/zeno -eval -numMsgToRecv 10 -metricsPipe /tmp/collect -client \
    -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr ${LISTEN_IP}:33000 -pkiLisAddr ${LISTEN_IP}:44000 \
    -pki ${PKI_IP}:33000 -pkiCertPath "/root/cert.pem" > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/
