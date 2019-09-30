#!/usr/bin/env bash

sleep 10

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575

# Gather some required context.
INTERNAL_IP=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip" -H "Metadata-Flavor: Google")
GCLOUD_SERVICE_ACC=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcloudServiceAcc" -H "Metadata-Flavor: Google")
GCLOUD_PROJECT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcloudProject" -H "Metadata-Flavor: Google")
GCLOUD_BUCKET=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/gcloudBucket" -H "Metadata-Flavor: Google")

# Pull operator binary from GCloud bucket.
/usr/bin/gsutil cp gs://acs-eval/operator /root/operator
chmod 0700 /root/operator

# Launch operator binary.
/root/operator -publicAddr 0.0.0.0:20443 -internalAddr ${INTERNAL_IP}:443 -gcloudServiceAcc ${GCLOUD_SERVICE_ACC} \
    -gcloudProject ${GCLOUD_PROJECT} -gcloudBucket ${GCLOUD_BUCKET} -certPath /root/operator-cert.pem -keyPath /root/operator-key.pem
