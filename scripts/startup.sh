#!/usr/bin/env bash

sleep 15

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575


# Retrieve metadata required for operation.

OPERATOR_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/operatorIP -H "Metadata-Flavor: Google")
EXP_ID=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/expID -H "Metadata-Flavor: Google")
WORKER_NAME=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode -H "Metadata-Flavor: Google")
EVAL_SYSTEM=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalSystem -H "Metadata-Flavor: Google")
NUM_CLIENTS=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/numClients -H "Metadata-Flavor: Google")
RESULT_FOLDER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/resultFolder -H "Metadata-Flavor: Google")

LISTEN_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull -H "Metadata-Flavor: Google")

PUNG_SERVER_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/pungServerIP -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig -H "Metadata-Flavor: Google")
KILL_ZENO_MIXES_IN_ROUND=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/killZenoMixesInRound -H "Metadata-Flavor: Google")

PUNG_CLIENTS_PER_PROC=$(( NUM_CLIENTS / 5))
printf "In case this is the Pung server machine, we will tell it to expect ${PUNG_CLIENTS_PER_PROC} messages per process.\n"


# Prepare to evaluate up to ten clients in case
# this is a clients machine. In case this is a
# server machine, values for CLIENT_01 will be
# exclusively used.

CLIENT_01=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client01 -H "Metadata-Flavor: Google")
CLIENT_01_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner01 -H "Metadata-Flavor: Google")
CLIENT_01_ADDR1="${LISTEN_IP}:33001"
CLIENT_01_ADDR2="${LISTEN_IP}:44001"
CLIENT_01_PATH="/root/${CLIENT_01}"
CLIENT_01_METRICS_PIPE=/tmp/collect01
CLIENT_01_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33001"
CLIENT_01_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}12"

CLIENT_02=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client02 -H "Metadata-Flavor: Google")
CLIENT_02_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner02 -H "Metadata-Flavor: Google")
CLIENT_02_ADDR1="${LISTEN_IP}:33002"
CLIENT_02_ADDR2="${LISTEN_IP}:44002"
CLIENT_02_PATH="/root/${CLIENT_02}"
CLIENT_02_METRICS_PIPE=/tmp/collect02
CLIENT_02_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33001"
CLIENT_02_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}12"

CLIENT_03=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client03 -H "Metadata-Flavor: Google")
CLIENT_03_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner03 -H "Metadata-Flavor: Google")
CLIENT_03_ADDR1="${LISTEN_IP}:33003"
CLIENT_03_ADDR2="${LISTEN_IP}:44003"
CLIENT_03_PATH="/root/${CLIENT_03}"
CLIENT_03_METRICS_PIPE=/tmp/collect03
CLIENT_03_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33002"
CLIENT_03_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}34"

CLIENT_04=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client04 -H "Metadata-Flavor: Google")
CLIENT_04_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner04 -H "Metadata-Flavor: Google")
CLIENT_04_ADDR1="${LISTEN_IP}:33004"
CLIENT_04_ADDR2="${LISTEN_IP}:44004"
CLIENT_04_PATH="/root/${CLIENT_04}"
CLIENT_04_METRICS_PIPE=/tmp/collect04
CLIENT_04_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33002"
CLIENT_04_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}34"

CLIENT_05=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client05 -H "Metadata-Flavor: Google")
CLIENT_05_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner05 -H "Metadata-Flavor: Google")
CLIENT_05_ADDR1="${LISTEN_IP}:33005"
CLIENT_05_ADDR2="${LISTEN_IP}:44005"
CLIENT_05_PATH="/root/${CLIENT_05}"
CLIENT_05_METRICS_PIPE=/tmp/collect05
CLIENT_05_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33003"
CLIENT_05_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}56"

CLIENT_06=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client06 -H "Metadata-Flavor: Google")
CLIENT_06_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner06 -H "Metadata-Flavor: Google")
CLIENT_06_ADDR1="${LISTEN_IP}:33006"
CLIENT_06_ADDR2="${LISTEN_IP}:44006"
CLIENT_06_PATH="/root/${CLIENT_06}"
CLIENT_06_METRICS_PIPE=/tmp/collect06
CLIENT_06_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33003"
CLIENT_06_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}56"

CLIENT_07=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client07 -H "Metadata-Flavor: Google")
CLIENT_07_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner07 -H "Metadata-Flavor: Google")
CLIENT_07_ADDR1="${LISTEN_IP}:33007"
CLIENT_07_ADDR2="${LISTEN_IP}:44007"
CLIENT_07_PATH="/root/${CLIENT_07}"
CLIENT_07_METRICS_PIPE=/tmp/collect07
CLIENT_07_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33004"
CLIENT_07_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}78"

CLIENT_08=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client08 -H "Metadata-Flavor: Google")
CLIENT_08_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner08 -H "Metadata-Flavor: Google")
CLIENT_08_ADDR1="${LISTEN_IP}:33008"
CLIENT_08_ADDR2="${LISTEN_IP}:44008"
CLIENT_08_PATH="/root/${CLIENT_08}"
CLIENT_08_METRICS_PIPE=/tmp/collect08
CLIENT_08_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33004"
CLIENT_08_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}78"

CLIENT_09=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client09 -H "Metadata-Flavor: Google")
CLIENT_09_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner09 -H "Metadata-Flavor: Google")
CLIENT_09_ADDR1="${LISTEN_IP}:33009"
CLIENT_09_ADDR2="${LISTEN_IP}:44009"
CLIENT_09_PATH="/root/${CLIENT_09}"
CLIENT_09_METRICS_PIPE=/tmp/collect09
CLIENT_09_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33005"
CLIENT_09_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}910"

CLIENT_10=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client10 -H "Metadata-Flavor: Google")
CLIENT_10_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner10 -H "Metadata-Flavor: Google")
CLIENT_10_ADDR1="${LISTEN_IP}:33010"
CLIENT_10_ADDR2="${LISTEN_IP}:44010"
CLIENT_10_PATH="/root/${CLIENT_10}"
CLIENT_10_METRICS_PIPE=/tmp/collect10
CLIENT_10_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33005"
CLIENT_10_PUNG_SHARED_SECRET="${WORKER_NAME}${LISTEN_IP}910"


# Prepare FIFO pipes for system and collector IPC.
mkfifo ${CLIENT_01_METRICS_PIPE} ${CLIENT_02_METRICS_PIPE} ${CLIENT_03_METRICS_PIPE} ${CLIENT_04_METRICS_PIPE} ${CLIENT_05_METRICS_PIPE} \
    ${CLIENT_06_METRICS_PIPE} ${CLIENT_07_METRICS_PIPE} ${CLIENT_08_METRICS_PIPE} ${CLIENT_09_METRICS_PIPE} ${CLIENT_10_METRICS_PIPE}
chmod 0600 ${CLIENT_01_METRICS_PIPE} ${CLIENT_02_METRICS_PIPE} ${CLIENT_03_METRICS_PIPE} ${CLIENT_04_METRICS_PIPE} ${CLIENT_05_METRICS_PIPE} \
    ${CLIENT_06_METRICS_PIPE} ${CLIENT_07_METRICS_PIPE} ${CLIENT_08_METRICS_PIPE} ${CLIENT_09_METRICS_PIPE} ${CLIENT_10_METRICS_PIPE}

# Prepare data paths for each client.
mkdir -p ${CLIENT_01_PATH} ${CLIENT_02_PATH} ${CLIENT_03_PATH} ${CLIENT_04_PATH} ${CLIENT_05_PATH} \
    ${CLIENT_06_PATH} ${CLIENT_07_PATH} ${CLIENT_08_PATH} ${CLIENT_09_PATH} ${CLIENT_10_PATH}


# Register with operator for current experiment.
curl --cacert /root/operator-cert.pem --request PUT --data-binary "{
    \"addresses\": [
        \"${CLIENT_01}\": \"${CLIENT_01_ADDR1}\",
        \"${CLIENT_02}\": \"${CLIENT_02_ADDR1}\",
        \"${CLIENT_03}\": \"${CLIENT_03_ADDR1}\",
        \"${CLIENT_04}\": \"${CLIENT_04_ADDR1}\",
        \"${CLIENT_05}\": \"${CLIENT_05_ADDR1}\",
        \"${CLIENT_06}\": \"${CLIENT_06_ADDR1}\",
        \"${CLIENT_07}\": \"${CLIENT_07_ADDR1}\",
        \"${CLIENT_08}\": \"${CLIENT_08_ADDR1}\",
        \"${CLIENT_09}\": \"${CLIENT_09_ADDR1}\",
        \"${CLIENT_10}\": \"${CLIENT_10_ADDR1}\"
    ]
}" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/register


# Pull files from GCloud bucket.
/usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/usr/bin/gsutil cp gs://acs-eval/collector /root/collector

if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then
    
    sleep 10

    printf "This is a Vuvuzela experiment, pull 'gs://acs-eval/vuvuzela-confs/pki.conf' as well\n"
    /usr/bin/gsutil cp gs://acs-eval/vuvuzela-confs/pki.conf /root/vuvuzela-confs/pki.conf

    while [ ! -e /root/vuvuzela-confs/pki.conf ]; do

        printf "Downloading 'gs://acs-eval/vuvuzela-confs/pki.conf' unsuccessful, trying again...\n"
        ls -lah /root/
        ls -lah /root/vuvuzela-confs/

        sleep 1

        /usr/bin/gsutil cp gs://acs-eval/vuvuzela-confs/pki.conf /root/vuvuzela-confs/pki.conf
    done
fi

tried=0
while ([ ! -e /root/${BINARY_TO_PULL} ] || [ ! -e /root/collector ]) && [ "${tried}" -lt 20 ]; do

    printf "Failed to pull required files from GCloud bucket, sleeping 1 second\n"
    ls -lah /root/

    sleep 1

    # Reattempt to pull files from GCloud bucket.
    /usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
    /usr/bin/gsutil cp gs://acs-eval/collector /root/collector

    tried=$(( tried + 1 ))
done

if [ "${tried}" -eq 20 ]; then

    printf "Waited 20 seconds for required experiment files to be downloaded, no success, shutting down\n"

    # Inform operator about failure to initialize.
    curl --cacert /root/operator-cert.pem --request PUT --data-binary "{
        \"failure\": \"waited 20 seconds for required experiment files to be downloaded from Storage bucket, no success, shutting down\"
    }" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/failed

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


if [ "${TYPE_OF_NODE}" == "server" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then

        # Add iptables rules to count network volume.
        iptables -A INPUT -p tcp --dport 33001
        iptables -A INPUT -p tcp --dport 44001
        iptables -A OUTPUT -p tcp --sport 33001
        iptables -A OUTPUT -p tcp --sport 44001
        iptables -Z -t filter -L INPUT
        iptables -Z -t filter -L OUTPUT

        # Run metrics collector sidecar in background.
        /root/collector -system zeno -server -pipe ${CLIENT_01_METRICS_PIPE} -metricsPath ${CLIENT_01_PATH}/ &

        # Run zeno as mix.
        /root/zeno -eval -killMixesInRound ${KILL_ZENO_MIXES_IN_ROUND} -metricsPipe ${CLIENT_01_METRICS_PIPE} -mix -name ${CLIENT_01} \
            -partner ${CLIENT_01_PARTNER} -msgPublicAddr ${CLIENT_01_ADDR1} -msgLisAddr ${CLIENT_01_ADDR1} -pkiLisAddr ${CLIENT_01_ADDR2} \
            -pki ${ZENO_PKI_IP}:33001 -pkiCertPath /root/operator-cert.pem > ${CLIENT_01_PATH}/log.evaluation
        
        # Wait for metrics collector to exit.
        wait

    else if [ "${EVAL_SYSTEM}" == "pung" ]; then

        iptables -A INPUT -p tcp --dport 33001
        iptables -A OUTPUT -p tcp --sport 33001
        iptables -Z -t filter -L INPUT
        iptables -Z -t filter -L OUTPUT

        /root/collector -system pung -server -pipe ${CLIENT_01_METRICS_PIPE} -metricsPath ${CLIENT_01_PATH}/ &

        /root/pung-server -e 30 -i ${LISTEN_IP} -s 33001 -n 5 -w 1 -p 0 -k 1 -t e -d 2 -b 0 -m ${PUNG_CLIENTS_PER_PROC} > ${CLIENT_01_PATH}/log.evaluation

        wait

    else if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then

        iptables -A INPUT -p tcp --dport 33001
        iptables -A OUTPUT -p tcp --sport 33001
        iptables -Z -t filter -L INPUT
        iptables -Z -t filter -L OUTPUT

        /root/collector -system vuvuzela -server -pipe ${CLIENT_01_METRICS_PIPE} -metricsPath ${CLIENT_01_PATH}/ &

        /root/vuvuzela-mix -eval -metricsPipe ${CLIENT_01_METRICS_PIPE} -addr ${CLIENT_01_ADDR1} -conf /root/vuvuzela-confs/${CLIENT_01}.conf \
            -pki /root/vuvuzela-confs/pki.conf > ${CLIENT_01_PATH}/log.evaluation

        wait

    fi

else if [ "${TYPE_OF_NODE}" == "coordinator" ]; then

    iptables -A INPUT -p tcp --dport 33001
    iptables -A OUTPUT -p tcp --sport 33001
    iptables -Z -t filter -L INPUT
    iptables -Z -t filter -L OUTPUT

    /root/collector -system vuvuzela -server -pipe ${CLIENT_01_METRICS_PIPE} -metricsPath ${CLIENT_01_PATH}/ &

    /root/vuvuzela-coordinator -eval -metricsPipe ${CLIENT_01_METRICS_PIPE} -addr ${ADDR1} \
        -wait 10s -pki /root/vuvuzela-confs/pki.conf > ${CLIENT_01_PATH}/log.evaluation

    wait

else if [ "${TYPE_OF_NODE}" == "client" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then



    else if [ "${EVAL_SYSTEM}" == "pung" ]; then



    else



    fi

fi


# Start all collectors and clients.

if [ "${CLIENT_01}" != "" ]; then
    CLIENT=${CLIENT_01} PARTNER=${CLIENT_01_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_01_ADDR1} ADDR2=${CLIENT_01_ADDR2} \
        CLIENT_PATH=${CLIENT_01_PATH} METRICS_PIPE=${CLIENT_01_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_01_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_01_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_02}" != "" ]; then
    CLIENT=${CLIENT_02} PARTNER=${CLIENT_02_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_02_ADDR1} ADDR2=${CLIENT_02_ADDR2} \
        CLIENT_PATH=${CLIENT_02_PATH} METRICS_PIPE=${CLIENT_02_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_02_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_02_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_03}" != "" ]; then
    CLIENT=${CLIENT_03} PARTNER=${CLIENT_03_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_03_ADDR1} ADDR2=${CLIENT_03_ADDR2} \
        CLIENT_PATH=${CLIENT_03_PATH} METRICS_PIPE=${CLIENT_03_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_03_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_03_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_04}" != "" ]; then
    CLIENT=${CLIENT_04} PARTNER=${CLIENT_04_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_04_ADDR1} ADDR2=${CLIENT_04_ADDR2} \
        CLIENT_PATH=${CLIENT_04_PATH} METRICS_PIPE=${CLIENT_04_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_04_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_04_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_05}" != "" ]; then
    CLIENT=${CLIENT_05} PARTNER=${CLIENT_05_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_05_ADDR1} ADDR2=${CLIENT_05_ADDR2} \
        CLIENT_PATH=${CLIENT_05_PATH} METRICS_PIPE=${CLIENT_05_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_05_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_05_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_06}" != "" ]; then
    CLIENT=${CLIENT_06} PARTNER=${CLIENT_06_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_06_ADDR1} ADDR2=${CLIENT_06_ADDR2} \
        CLIENT_PATH=${CLIENT_06_PATH} METRICS_PIPE=${CLIENT_06_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_06_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_06_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_07}" != "" ]; then
    CLIENT=${CLIENT_07} PARTNER=${CLIENT_07_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_07_ADDR1} ADDR2=${CLIENT_07_ADDR2} \
        CLIENT_PATH=${CLIENT_07_PATH} METRICS_PIPE=${CLIENT_07_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_07_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_07_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_08}" != "" ]; then
    CLIENT=${CLIENT_08} PARTNER=${CLIENT_08_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_08_ADDR1} ADDR2=${CLIENT_08_ADDR2} \
        CLIENT_PATH=${CLIENT_08_PATH} METRICS_PIPE=${CLIENT_08_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_08_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_08_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_09}" != "" ]; then
    CLIENT=${CLIENT_09} PARTNER=${CLIENT_09_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_09_ADDR1} ADDR2=${CLIENT_09_ADDR2} \
        CLIENT_PATH=${CLIENT_09_PATH} METRICS_PIPE=${CLIENT_09_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_09_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_09_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi

if [ "${CLIENT_10}" != "" ]; then
    CLIENT=${CLIENT_10} PARTNER=${CLIENT_10_PARTNER} LISTEN_IP=${LISTEN_IP} ADDR1=${CLIENT_10_ADDR1} ADDR2=${CLIENT_10_ADDR2} \
        CLIENT_PATH=${CLIENT_10_PATH} METRICS_PIPE=${CLIENT_10_METRICS_PIPE} ZENO_PKI_IP=${OPERATOR_IP} PUNG_SERVER_ADDR=${CLIENT_10_PUNG_SERVER_ADDR} \
        PUNG_CLIENTS_PER_PROC=${PUNG_CLIENTS_PER_PROC} KILL_ZENO_MIXES_IN_ROUND=${KILL_ZENO_MIXES_IN_ROUND} PUNG_SHARED_SECRET=${CLIENT_10_PUNG_SHARED_SECRET} \
        /bin/bash /root/${EVAL_SCRIPT_TO_PULL}
fi


# Reset tc configuration.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc del dev ${NET_DEVICE} root
fi


# Upload result files to GCloud bucket.

if [ "${CLIENT_01}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_01_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_01}/
fi

if [ "${CLIENT_02}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_02_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_02}/
fi

if [ "${CLIENT_03}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_03_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_03}/
fi

if [ "${CLIENT_04}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_04_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_04}/
fi

if [ "${CLIENT_05}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_05_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_05}/
fi

if [ "${CLIENT_06}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_06_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_06}/
fi

if [ "${CLIENT_07}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_07_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_07}/
fi

if [ "${CLIENT_08}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_08_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_08}/
fi

if [ "${CLIENT_09}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_09_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_09}/
fi

if [ "${CLIENT_10}" != "" ]; then
    /usr/bin/gsutil -m cp ${CLIENT_10_PATH}/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${WORKER_NAME}_${LISTEN_IP}/${CLIENT_10}/
fi


# Mark worker as finished at operator.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${WORKER_NAME}/finished
