#!/usr/bin/env bash

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki.pem ~/cert.pem
chmod 0644 ~/cert.pem

# Configure tc according to environment variable.

# Run zeno as client.
~/zeno -client -msgPublicAddr ${EXTERNAL_IP}:33000 -msgLisAddr 0.0.0.0:33000 -pkiLisAddr 0.0.0.0:44000 -pki ${PKI_IP}:33000 -pkiCertPath "~/cert.pem"

# Reset tc configuration.

# Upload result files to GCloud bucket.
