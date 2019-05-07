#!/usr/bin/env bash

upload_cert() {

    # Wait some time to allow PKI to bootstrap.
    sleep 15

    # Upload TLS certificate to bucket.
    /snap/bin/gsutil cp ./cert.pem gs://acs-eval/cert_zeno-pki.pem
}

# Start process in background to eventually
# upload the PKI process' TLS certificate.
upload_cert &

# Run main PKI process.
~/zeno-pki -publicAddr ${LISTEN_IP}:33000 -listenAddr 0.0.0.0:33000 -controlPlaneAddr 0.0.0.0:26345
