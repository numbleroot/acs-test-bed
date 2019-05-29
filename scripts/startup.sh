#!/usr/bin/env bash

sleep 5

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575

# Prepare FIFO pipe for system and collector IPC.
mkfifo /tmp/collect
chmod 0600 /tmp/collect

# Retrieve metadata required for operation.
LISTEN_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip -H "Metadata-Flavor: Google")
NAME_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode -H "Metadata-Flavor: Google")
PARTNER_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partnerOfNode -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode -H "Metadata-Flavor: Google")
RESULT_FOLDER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/resultFolder -H "Metadata-Flavor: Google")
EVAL_SCRIPT_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalScriptToPull -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig -H "Metadata-Flavor: Google")
KILL_ZENO_MIXES_IN_ROUND=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/killZenoMixesInRound -H "Metadata-Flavor: Google")
PKI_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/pkiIP -H "Metadata-Flavor: Google")

tried=0
while [ ! -e /snap/bin/gsutil ] && [ "${tried}" -lt 30 ]; do

    printf "/snap/bin/gsutil not yet available, sleeping 1 second\n"
    ls -lah /snap/bin/

    sleep 1
    tried=$(( tried + 1 ))
done

if [ "${tried}" -eq 30 ]; then
    printf "Waited 30 seconds for /snap/bin/gsutil to become available, no success, shutting down\n"
    poweroff
fi

tried=0
while [ ! -e /usr/bin/python2 ] && [ "${tried}" -lt 30 ]; do

    printf "/usr/bin/python2 not yet available, sleeping 1 second\n"
    ls -lah /usr/bin | grep pyth

    sleep 1
    tried=$(( tried + 1 ))
done

if [ "${tried}" -eq 30 ]; then
    printf "Waited 30 seconds for /usr/bin/python2 to become available, no success, shutting down\n"
    poweroff
fi

sleep 10

# Pull files from GCloud bucket.
/snap/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} /root/${EVAL_SCRIPT_TO_PULL}
/snap/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/snap/bin/gsutil cp gs://acs-eval/collector /root/collector

tried=0
while ([ ! -e /root/${EVAL_SCRIPT_TO_PULL} ] || [ ! -e /root/${BINARY_TO_PULL} ] || [ ! -e /root/collector ]) && [ "${tried}" -lt 20 ]; do

    printf "Failed to pull required files from GCloud bucket, sleeping 1 second\n"
    ls -lah /root/

    sleep 1

    # Reattempt to pull files from GCloud bucket.
    /snap/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} /root/${EVAL_SCRIPT_TO_PULL}
    /snap/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
    /snap/bin/gsutil cp gs://acs-eval/collector /root/collector
done

if [ "${tried}" -eq 20 ]; then
    printf "Waited 20 seconds for required experiment files to be downloaded, no success, shutting down\n"
    poweroff
fi

# Make the downloaded binaries executable.
chmod 0700 /root/${BINARY_TO_PULL}
chmod 0700 /root/collector

# Hand over to evaluation script.
LISTEN_IP=${LISTEN_IP} NAME_OF_NODE=${NAME_OF_NODE} PARTNER_OF_NODE=${PARTNER_OF_NODE} \
    TYPE_OF_NODE=${TYPE_OF_NODE} RESULT_FOLDER=${RESULT_FOLDER} TC_CONFIG=${TC_CONFIG} \
    KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PKI_IP=${PKI_IP} /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
