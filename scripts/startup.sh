#!/usr/bin/env bash

sleep 7

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575

# Prepare FIFO pipe for system and collector IPC.
mkfifo /tmp/collect
chmod 0600 /tmp/collect

# Add iptables rule to be able to count number of transferred
# bytes over evaluation system port.
iptables -A INPUT -p tcp --dport 33000
iptables -A OUTPUT -p tcp --dport 33000

# Retrieve metadata required for operation.
LISTEN_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip -H "Metadata-Flavor: Google")
NAME_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode -H "Metadata-Flavor: Google")
RESULT_FOLDER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/resultFolder -H "Metadata-Flavor: Google")
EVAL_SCRIPT_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalScriptToPull -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig -H "Metadata-Flavor: Google")
PKI_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/pkiIP -H "Metadata-Flavor: Google")

# Pull files from GCloud bucket.
/snap/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} /root/${EVAL_SCRIPT_TO_PULL}
/snap/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/snap/bin/gsutil cp gs://acs-eval/collector /root/collector

# Make the downloaded binaries executable.
chmod 0700 /root/${BINARY_TO_PULL}
chmod 0700 /root/collector

# Hand over to evaluation script.
LISTEN_IP=${LISTEN_IP} NAME_OF_NODE=${NAME_OF_NODE} TYPE_OF_NODE=${TYPE_OF_NODE} RESULT_FOLDER=${RESULT_FOLDER} TC_CONFIG=${TC_CONFIG} PKI_IP=${PKI_IP} /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
