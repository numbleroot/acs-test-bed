#!/usr/bin/env bash

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki.pem /root/cert.pem
chmod 0644 /root/cert.pem

# Configure tc according to environment variable.

# Run zeno as client.
/root/zeno -eval -numMsgToRecv 10 -client -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr ${LISTEN_IP}:33000 -pkiLisAddr ${LISTEN_IP}:44000 -pki ${PKI_IP}:33000 -pkiCertPath "/root/cert.pem" > /root/zeno_client_${LISTEN_IP}.log

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil cp /root/zeno_client_${LISTEN_IP}.log gs://acs-eval/${RESULT_FOLDER}/zeno_client_${LISTEN_IP}.log
