#!/usr/bin/env bash

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki.pem /root/cert.pem
chmod 0644 /root/cert.pem

# Configure tc according to environment variable.

# Run metrics collector sidecar in background.
# /root/collector -client -pipe /tmp/collect -metricsPath /root/ &

# Run zeno as mix.
/root/zeno -mix -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr ${LISTEN_IP}:33000 -pkiLisAddr ${LISTEN_IP}:44000 -pki ${PKI_IP}:33000 -pkiCertPath "/root/cert.pem" > /root/zeno_mix_${LISTEN_IP}_log.evaluation

# Wait for metrics collector to exit.
# wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/${LISTEN_IP}/
