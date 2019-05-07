#!/usr/bin/env bash

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki.pem ~/cert.pem
chmod 0644 ~/cert.pem

# Configure tc according to environment variable.

# Run zeno as client.
~/zeno -eval -numMsgToRecv 10 -client -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr 0.0.0.0:33000 -pkiLisAddr 0.0.0.0:44000 -pki ${PKI_IP}:33000 -pkiCertPath "~/cert.pem" > ~/zeno_client_${LISTEN_IP}.log

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil cp ~/zeno_client_${LISTEN_IP}.log gs://acs-eval/${RESULT_FOLDER}/zeno_client_${LISTEN_IP}.log
