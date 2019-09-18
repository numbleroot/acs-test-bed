#!/usr/bin/env bash

sleep 15

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575


# Retrieve metadata required for operation.
LISTEN_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip -H "Metadata-Flavor: Google")
OPERATOR_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/operatorIP -H "Metadata-Flavor: Google")
EXP_ID=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/expID -H "Metadata-Flavor: Google")
PUNG_SERVER_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/pungServerIP -H "Metadata-Flavor: Google")
WORKER_NAME=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode -H "Metadata-Flavor: Google")
RESULT_FOLDER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/resultFolder -H "Metadata-Flavor: Google")
EVAL_SCRIPT_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalScriptToPull -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig -H "Metadata-Flavor: Google")
KILL_ZENO_MIXES_IN_ROUND=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/killZenoMixesInRound -H "Metadata-Flavor: Google")

# Prepare to evaluate eight clients in case
# this is a clients machine.
CLIENT_1=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client1 -H "Metadata-Flavor: Google")
PARTNER_1=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner1 -H "Metadata-Flavor: Google")
CLIENT_2=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client2 -H "Metadata-Flavor: Google")
PARTNER_2=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner2 -H "Metadata-Flavor: Google")
CLIENT_3=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client3 -H "Metadata-Flavor: Google")
PARTNER_3=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner3 -H "Metadata-Flavor: Google")
CLIENT_4=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client4 -H "Metadata-Flavor: Google")
PARTNER_4=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner4 -H "Metadata-Flavor: Google")
CLIENT_5=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client5 -H "Metadata-Flavor: Google")
PARTNER_5=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner5 -H "Metadata-Flavor: Google")
CLIENT_6=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client6 -H "Metadata-Flavor: Google")
PARTNER_6=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner6 -H "Metadata-Flavor: Google")
CLIENT_7=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client7 -H "Metadata-Flavor: Google")
PARTNER_7=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner7 -H "Metadata-Flavor: Google")
CLIENT_8=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client8 -H "Metadata-Flavor: Google")
PARTNER_8=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner8 -H "Metadata-Flavor: Google")


# Prepare FIFO pipes for system and collector IPC.
mkfifo /tmp/collect1 /tmp/collect2 /tmp/collect3 /tmp/collect4 /tmp/collect5 /tmp/collect6 /tmp/collect7 /tmp/collect8
chmod 0600 /tmp/collect1 /tmp/collect2 /tmp/collect3 /tmp/collect4 /tmp/collect5 /tmp/collect6 /tmp/collect7 /tmp/collect8

# Prepare metric paths for each client.
mkdir -p /root/${CLIENT_1} /root/${CLIENT_2} /root/${CLIENT_3} /root/${CLIENT_4} /root/${CLIENT_5} /root/${CLIENT_6} /root/${CLIENT_7} /root/${CLIENT_8}


# Register with operator for current experiment.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/register


# Pull files from GCloud bucket.
/usr/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} /root/${EVAL_SCRIPT_TO_PULL}
/usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/usr/bin/gsutil cp gs://acs-eval/collector /root/collector

tried=0
while ([ ! -e /root/${EVAL_SCRIPT_TO_PULL} ] || [ ! -e /root/${BINARY_TO_PULL} ] || [ ! -e /root/collector ]) && [ "${tried}" -lt 20 ]; do

    printf "Failed to pull required files from GCloud bucket, sleeping 1 second\n"
    ls -lah /root/

    sleep 1

    # Reattempt to pull files from GCloud bucket.
    /usr/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} /root/${EVAL_SCRIPT_TO_PULL}
    /usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
    /usr/bin/gsutil cp gs://acs-eval/collector /root/collector
done

if [ "${tried}" -eq 20 ]; then

    printf "Waited 20 seconds for required experiment files to be downloaded, no success, shutting down\n"

    # Inform operator about failure to initialize.
    curl --cacert /root/operator-cert.pem --request PUT --data-binary "{
        \"Failure\": \"Waited 20 seconds for required experiment files to be downloaded from Storage bucket, no success, shutting down\"
    }" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/failure

    poweroff
fi

# Make the downloaded binaries executable.
chmod 0700 /root/${BINARY_TO_PULL}
chmod 0700 /root/collector

sleep 5

# Signal readiness of process to experiment script.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/ready


# Determine active network device.
NET_DEVICE=$(ip addr | awk '/state UP/ {print $2}' | sed 's/.$//')

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


# Start all collectors and clients.

if [ "${CLIENT_1}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_1} PARTNER=${PARTNER_1} METRICS_PIPE=/tmp/collect1 \
        CLIENT_PATH=/root/${CLIENT_1} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}12 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_2}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_2} PARTNER=${PARTNER_2} METRICS_PIPE=/tmp/collect2 \
        CLIENT_PATH=/root/${CLIENT_2} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}12 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_3}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_3} PARTNER=${PARTNER_3} METRICS_PIPE=/tmp/collect3 \
        CLIENT_PATH=/root/${CLIENT_3} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}34 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_4}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_4} PARTNER=${PARTNER_4} METRICS_PIPE=/tmp/collect4 \
        CLIENT_PATH=/root/${CLIENT_4} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}34 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_5}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_5} PARTNER=${PARTNER_5} METRICS_PIPE=/tmp/collect5 \
        CLIENT_PATH=/root/${CLIENT_5} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}56 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_6}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_6} PARTNER=${PARTNER_6} METRICS_PIPE=/tmp/collect6 \
        CLIENT_PATH=/root/${CLIENT_6} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}56 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_7}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_7} PARTNER=${PARTNER_7} METRICS_PIPE=/tmp/collect7 \
        CLIENT_PATH=/root/${CLIENT_7} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}78 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_8}" != "" ]; then
    LISTEN_IP=${LISTEN_IP} CLIENT=${CLIENT_8} PARTNER=${PARTNER_8} METRICS_PIPE=/tmp/collect8 \
        CLIENT_PATH=/root/${CLIENT_8} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_IP=${PUNG_SERVER_IP} \
        KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${WORKER_NAME}${LISTEN_IP}78 \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi


# Reset tc configuration.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc del dev ${NET_DEVICE} root
fi


# Upload result files to GCloud bucket.

if [ "${CLIENT_1}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_1}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_1}/
fi

if [ "${CLIENT_2}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_2}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_2}/
fi

if [ "${CLIENT_3}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_3}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_3}/
fi

if [ "${CLIENT_4}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_4}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_4}/
fi

if [ "${CLIENT_5}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_5}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_5}/
fi

if [ "${CLIENT_6}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_6}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_6}/
fi

if [ "${CLIENT_7}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_7}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_7}/
fi

if [ "${CLIENT_8}" != "" ]; then
    /usr/bin/gsutil -m cp /root/${CLIENT_8}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_8}/
fi


# Mark worker as finished at operator.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/finished
