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
NAME_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode -H "Metadata-Flavor: Google")
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
CLIENT_01_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33001"
CLIENT_01_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}12"

CLIENT_02=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client02 -H "Metadata-Flavor: Google")
CLIENT_02_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner02 -H "Metadata-Flavor: Google")
CLIENT_02_ADDR1="${LISTEN_IP}:33002"
CLIENT_02_ADDR2="${LISTEN_IP}:44002"
CLIENT_02_PATH="/root/${CLIENT_02}"
CLIENT_02_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33001"
CLIENT_02_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}12"

CLIENT_03=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client03 -H "Metadata-Flavor: Google")
CLIENT_03_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner03 -H "Metadata-Flavor: Google")
CLIENT_03_ADDR1="${LISTEN_IP}:33003"
CLIENT_03_ADDR2="${LISTEN_IP}:44003"
CLIENT_03_PATH="/root/${CLIENT_03}"
CLIENT_03_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33002"
CLIENT_03_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}34"

CLIENT_04=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client04 -H "Metadata-Flavor: Google")
CLIENT_04_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner04 -H "Metadata-Flavor: Google")
CLIENT_04_ADDR1="${LISTEN_IP}:33004"
CLIENT_04_ADDR2="${LISTEN_IP}:44004"
CLIENT_04_PATH="/root/${CLIENT_04}"
CLIENT_04_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33002"
CLIENT_04_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}34"

CLIENT_05=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client05 -H "Metadata-Flavor: Google")
CLIENT_05_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner05 -H "Metadata-Flavor: Google")
CLIENT_05_ADDR1="${LISTEN_IP}:33005"
CLIENT_05_ADDR2="${LISTEN_IP}:44005"
CLIENT_05_PATH="/root/${CLIENT_05}"
CLIENT_05_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33003"
CLIENT_05_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}56"

CLIENT_06=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client06 -H "Metadata-Flavor: Google")
CLIENT_06_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner06 -H "Metadata-Flavor: Google")
CLIENT_06_ADDR1="${LISTEN_IP}:33006"
CLIENT_06_ADDR2="${LISTEN_IP}:44006"
CLIENT_06_PATH="/root/${CLIENT_06}"
CLIENT_06_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33003"
CLIENT_06_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}56"

CLIENT_07=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client07 -H "Metadata-Flavor: Google")
CLIENT_07_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner07 -H "Metadata-Flavor: Google")
CLIENT_07_ADDR1="${LISTEN_IP}:33007"
CLIENT_07_ADDR2="${LISTEN_IP}:44007"
CLIENT_07_PATH="/root/${CLIENT_07}"
CLIENT_07_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33004"
CLIENT_07_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}78"

CLIENT_08=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client08 -H "Metadata-Flavor: Google")
CLIENT_08_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner08 -H "Metadata-Flavor: Google")
CLIENT_08_ADDR1="${LISTEN_IP}:33008"
CLIENT_08_ADDR2="${LISTEN_IP}:44008"
CLIENT_08_PATH="/root/${CLIENT_08}"
CLIENT_08_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33004"
CLIENT_08_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}78"

CLIENT_09=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client09 -H "Metadata-Flavor: Google")
CLIENT_09_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner09 -H "Metadata-Flavor: Google")
CLIENT_09_ADDR1="${LISTEN_IP}:33009"
CLIENT_09_ADDR2="${LISTEN_IP}:44009"
CLIENT_09_PATH="/root/${CLIENT_09}"
CLIENT_09_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33005"
CLIENT_09_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}910"

CLIENT_10=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/client10 -H "Metadata-Flavor: Google")
CLIENT_10_PARTNER=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner10 -H "Metadata-Flavor: Google")
CLIENT_10_ADDR1="${LISTEN_IP}:33010"
CLIENT_10_ADDR2="${LISTEN_IP}:44010"
CLIENT_10_PATH="/root/${CLIENT_10}"
CLIENT_10_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33005"
CLIENT_10_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}910"


# Prepare FIFO pipe for system and collector IPC.
mkfifo /tmp/collect
chmod 0600 /tmp/collect

# Prepare data paths for each client.
mkdir -p ${CLIENT_01_PATH} ${CLIENT_02_PATH} ${CLIENT_03_PATH} ${CLIENT_04_PATH} ${CLIENT_05_PATH} \
    ${CLIENT_06_PATH} ${CLIENT_07_PATH} ${CLIENT_08_PATH} ${CLIENT_09_PATH} ${CLIENT_10_PATH}


# Register with operator for current experiment.
curl --cacert /root/operator-cert.pem --request PUT --data-binary "{
    \"address\": \"${CLIENT_01_ADDR1}\"
}" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/register


# Pull files from GCloud bucket.
/usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/usr/bin/gsutil cp gs://acs-eval/collector /root/collector

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
    }" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/failed

    poweroff
fi

if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then
    
    sleep 10

    printf "This is a Vuvuzela experiment, pull 'gs://acs-eval/vuvuzela-confs/pki.conf' as well\n"
    /usr/bin/gsutil cp gs://acs-eval/vuvuzela-confs/pki.conf /root/vuvuzela-confs/pki.conf

    while [ ! -e /root/vuvuzela-confs/pki.conf ]; do

        printf "Download of 'gs://acs-eval/vuvuzela-confs/pki.conf' unsuccessful, trying again...\n"
        ls -lah /root/
        ls -lah /root/vuvuzela-confs/

        sleep 1

        /usr/bin/gsutil cp gs://acs-eval/vuvuzela-confs/pki.conf /root/vuvuzela-confs/pki.conf
    done
fi

# Make the downloaded binaries executable.
chmod 0700 /root/${BINARY_TO_PULL}
chmod 0700 /root/collector


# Prepare some surroundings logging.
echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_01_PATH}/log.evaluation
echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_01_PATH}/log.evaluation
echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_01_PATH}/log.evaluation
echo "System info: '$(uname -a)'\n" >> ${CLIENT_01_PATH}/log.evaluation
echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_01_PATH}/log.evaluation
echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_01_PATH}/log.evaluation
echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_01_PATH}/log.evaluation

if [ "${TYPE_OF_NODE}" == "client" ]; then

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_02_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_02_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_02_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_02_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_02_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_02_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_02_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_03_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_03_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_03_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_03_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_03_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_03_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_03_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_04_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_04_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_04_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_04_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_04_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_04_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_04_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_05_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_05_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_05_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_05_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_05_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_05_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_05_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_06_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_06_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_06_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_06_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_06_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_06_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_06_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_07_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_07_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_07_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_07_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_07_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_07_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_07_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_08_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_08_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_08_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_08_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_08_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_08_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_08_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_09_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_09_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_09_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_09_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_09_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_09_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_09_PATH}/log.evaluation

    echo "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > ${CLIENT_10_PATH}/log.evaluation
    echo "Result folder: '${RESULT_FOLDER}'\n" >> ${CLIENT_10_PATH}/log.evaluation
    echo "${NUM_CLIENTS} clients will participate, TC parameters set to: '${TC_CONFIG}'.\n" >> ${CLIENT_10_PATH}/log.evaluation
    echo "System info: '$(uname -a)'\n" >> ${CLIENT_10_PATH}/log.evaluation
    echo "CPU: $(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')\n" >> ${CLIENT_10_PATH}/log.evaluation
    echo "Memory: $(lsmem | grep "Total online memory" | awk '{print $4}')\n" >> ${CLIENT_10_PATH}/log.evaluation
    echo "Storage: $(lsblk -o TYPE,SIZE,MODEL | grep disk | awk '{print $2,$3}')\n" >> ${CLIENT_10_PATH}/log.evaluation

fi


sleep 5

# Signal readiness of process to experiment script.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/ready


# Determine active network device.
NET_DEVICE=$(ip addr | awk '/state UP/ {print $2}' | sed 's/.$//')

# Configure tc according to environment variable.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc add dev ${NET_DEVICE} root ${TC_CONFIG}
    printf "Configured ${NET_DEVICE} with tc parameters.\n"
fi


# Add iptables rules to count network volume.

iptables -t filter -A INPUT -p tcp --sport 33001
iptables -t filter -A INPUT -p tcp --dport 33001
iptables -t filter -A OUTPUT -p tcp --sport 33001
iptables -t filter -A OUTPUT -p tcp --dport 33001

iptables -t filter -A INPUT -p tcp --sport 33002
iptables -t filter -A INPUT -p tcp --dport 33002
iptables -t filter -A OUTPUT -p tcp --sport 33002
iptables -t filter -A OUTPUT -p tcp --dport 33002

iptables -t filter -A INPUT -p tcp --sport 33003
iptables -t filter -A INPUT -p tcp --dport 33003
iptables -t filter -A OUTPUT -p tcp --sport 33003
iptables -t filter -A OUTPUT -p tcp --dport 33003

iptables -t filter -A INPUT -p tcp --sport 33004
iptables -t filter -A INPUT -p tcp --dport 33004
iptables -t filter -A OUTPUT -p tcp --sport 33004
iptables -t filter -A OUTPUT -p tcp --dport 33004

iptables -t filter -A INPUT -p tcp --sport 33005
iptables -t filter -A INPUT -p tcp --dport 33005
iptables -t filter -A OUTPUT -p tcp --sport 33005
iptables -t filter -A OUTPUT -p tcp --dport 33005

iptables -t filter -A INPUT -p tcp --sport 33006
iptables -t filter -A INPUT -p tcp --dport 33006
iptables -t filter -A OUTPUT -p tcp --sport 33006
iptables -t filter -A OUTPUT -p tcp --dport 33006

iptables -t filter -A INPUT -p tcp --sport 33007
iptables -t filter -A INPUT -p tcp --dport 33007
iptables -t filter -A OUTPUT -p tcp --sport 33007
iptables -t filter -A OUTPUT -p tcp --dport 33007

iptables -t filter -A INPUT -p tcp --sport 33008
iptables -t filter -A INPUT -p tcp --dport 33008
iptables -t filter -A OUTPUT -p tcp --sport 33008
iptables -t filter -A OUTPUT -p tcp --dport 33008

iptables -t filter -A INPUT -p tcp --sport 33009
iptables -t filter -A INPUT -p tcp --dport 33009
iptables -t filter -A OUTPUT -p tcp --sport 33009
iptables -t filter -A OUTPUT -p tcp --dport 33009

iptables -t filter -A INPUT -p tcp --sport 33010
iptables -t filter -A INPUT -p tcp --dport 33010
iptables -t filter -A OUTPUT -p tcp --sport 33010
iptables -t filter -A OUTPUT -p tcp --dport 33010

iptables -t filter -A INPUT -p tcp --sport 44001
iptables -t filter -A INPUT -p tcp --dport 44001
iptables -t filter -A OUTPUT -p tcp --sport 44001
iptables -t filter -A OUTPUT -p tcp --dport 44001

iptables -t filter -A INPUT -p tcp --sport 44002
iptables -t filter -A INPUT -p tcp --dport 44002
iptables -t filter -A OUTPUT -p tcp --sport 44002
iptables -t filter -A OUTPUT -p tcp --dport 44002

iptables -t filter -A INPUT -p tcp --sport 44003
iptables -t filter -A INPUT -p tcp --dport 44003
iptables -t filter -A OUTPUT -p tcp --sport 44003
iptables -t filter -A OUTPUT -p tcp --dport 44003

iptables -t filter -A INPUT -p tcp --sport 44004
iptables -t filter -A INPUT -p tcp --dport 44004
iptables -t filter -A OUTPUT -p tcp --sport 44004
iptables -t filter -A OUTPUT -p tcp --dport 44004

iptables -t filter -A INPUT -p tcp --sport 44005
iptables -t filter -A INPUT -p tcp --dport 44005
iptables -t filter -A OUTPUT -p tcp --sport 44005
iptables -t filter -A OUTPUT -p tcp --dport 44005

iptables -t filter -A INPUT -p tcp --sport 44006
iptables -t filter -A INPUT -p tcp --dport 44006
iptables -t filter -A OUTPUT -p tcp --sport 44006
iptables -t filter -A OUTPUT -p tcp --dport 44006

iptables -t filter -A INPUT -p tcp --sport 44007
iptables -t filter -A INPUT -p tcp --dport 44007
iptables -t filter -A OUTPUT -p tcp --sport 44007
iptables -t filter -A OUTPUT -p tcp --dport 44007

iptables -t filter -A INPUT -p tcp --sport 44008
iptables -t filter -A INPUT -p tcp --dport 44008
iptables -t filter -A OUTPUT -p tcp --sport 44008
iptables -t filter -A OUTPUT -p tcp --dport 44008

iptables -t filter -A INPUT -p tcp --sport 44009
iptables -t filter -A INPUT -p tcp --dport 44009
iptables -t filter -A OUTPUT -p tcp --sport 44009
iptables -t filter -A OUTPUT -p tcp --dport 44009

iptables -t filter -A INPUT -p tcp --sport 44010
iptables -t filter -A INPUT -p tcp --dport 44010
iptables -t filter -A OUTPUT -p tcp --sport 44010
iptables -t filter -A OUTPUT -p tcp --dport 44010

iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT


# Run metrics collector sidecar in background.
/root/collector -system ${EVAL_SYSTEM} -typeOfNode ${TYPE_OF_NODE} -pipe /tmp/collect -metricsPath /root/ &


if [ "${TYPE_OF_NODE}" == "server" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_01_PATH}/log.evaluation

        # Run zeno as mix.
        /root/zeno -eval -killMixesInRound ${KILL_ZENO_MIXES_IN_ROUND} -metricsPipe /tmp/collect -mix -name ${CLIENT_01} \
            -partner ${CLIENT_01_PARTNER} -msgPublicAddr ${CLIENT_01_ADDR1} -msgLisAddr ${CLIENT_01_ADDR1} -pkiLisAddr ${CLIENT_01_ADDR2} \
            -pki ${OPERATOR_IP}:44001 -pkiCertPath /root/operator-cert.pem >> ${CLIENT_01_PATH}/log.evaluation

    else if [ "${EVAL_SYSTEM}" == "pung" ]; then

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_01_PATH}/log.evaluation

        # Run Pung's server.
        /root/pung-server -e 30 -i ${LISTEN_IP} -s 33001 -n 5 -w 1 -p 0 -k 1 -t e -d 2 -b 0 -m ${PUNG_CLIENTS_PER_PROC} >> ${CLIENT_01_PATH}/log.evaluation

    else if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then

        echo "\n" >> ${CLIENT_01_PATH}/log.evaluation

        # Run mix component of Vuvuzela.
        /root/vuvuzela-mix -eval -metricsPipe /tmp/collect -addr ${CLIENT_01_ADDR1} -conf /root/vuvuzela-confs/${CLIENT_01}.conf \
            -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_01_PATH}/log.evaluation

    fi

else if [ "${TYPE_OF_NODE}" == "coordinator" ]; then

    echo "\n" >> ${CLIENT_01_PATH}/log.evaluation

    # Run coordinator component of Vuvuzela.
    /root/vuvuzela-coordinator -eval -metricsPipe /tmp/collect -addr ${CLIENT_01_ADDR1} \
        -wait 10s -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_01_PATH}/log.evaluation

else if [ "${TYPE_OF_NODE}" == "client" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then

        # Run ten zeno clients.

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_01_PATH}/log.evaluation        
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_01} -partner ${CLIENT_01_PARTNER} \
            -msgPublicAddr ${CLIENT_01_ADDR1} -msgLisAddr ${CLIENT_01_ADDR1} -pkiLisAddr ${CLIENT_01_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_01_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_02_PATH}/log.evaluation 
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_02} -partner ${CLIENT_02_PARTNER} \
            -msgPublicAddr ${CLIENT_02_ADDR1} -msgLisAddr ${CLIENT_02_ADDR1} -pkiLisAddr ${CLIENT_02_ADDR2} -pki ${OPERATOR_IP}:44002 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_02_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_03_PATH}/log.evaluation 
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_03} -partner ${CLIENT_03_PARTNER} \
            -msgPublicAddr ${CLIENT_03_ADDR1} -msgLisAddr ${CLIENT_03_ADDR1} -pkiLisAddr ${CLIENT_03_ADDR2} -pki ${OPERATOR_IP}:44003 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_03_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_04_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_04} -partner ${CLIENT_04_PARTNER} \
            -msgPublicAddr ${CLIENT_04_ADDR1} -msgLisAddr ${CLIENT_04_ADDR1} -pkiLisAddr ${CLIENT_04_ADDR2} -pki ${OPERATOR_IP}:44004 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_04_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_05_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_05} -partner ${CLIENT_05_PARTNER} \
            -msgPublicAddr ${CLIENT_05_ADDR1} -msgLisAddr ${CLIENT_05_ADDR1} -pkiLisAddr ${CLIENT_05_ADDR2} -pki ${OPERATOR_IP}:44005 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_05_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_06_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_06} -partner ${CLIENT_06_PARTNER} \
            -msgPublicAddr ${CLIENT_06_ADDR1} -msgLisAddr ${CLIENT_06_ADDR1} -pkiLisAddr ${CLIENT_06_ADDR2} -pki ${OPERATOR_IP}:44006 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_06_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_07_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_07} -partner ${CLIENT_07_PARTNER} \
            -msgPublicAddr ${CLIENT_07_ADDR1} -msgLisAddr ${CLIENT_07_ADDR1} -pkiLisAddr ${CLIENT_07_ADDR2} -pki ${OPERATOR_IP}:44007 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_07_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_08_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_08} -partner ${CLIENT_08_PARTNER} \
            -msgPublicAddr ${CLIENT_08_ADDR1} -msgLisAddr ${CLIENT_08_ADDR1} -pkiLisAddr ${CLIENT_08_ADDR2} -pki ${OPERATOR_IP}:44008 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_08_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_09_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_09} -partner ${CLIENT_09_PARTNER} \
            -msgPublicAddr ${CLIENT_09_ADDR1} -msgLisAddr ${CLIENT_09_ADDR1} -pkiLisAddr ${CLIENT_09_ADDR2} -pki ${OPERATOR_IP}:44009 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_09_PATH}/log.evaluation

        echo "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> ${CLIENT_10_PATH}/log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -client -name ${CLIENT_10} -partner ${CLIENT_10_PARTNER} \
            -msgPublicAddr ${CLIENT_10_ADDR1} -msgLisAddr ${CLIENT_10_ADDR1} -pkiLisAddr ${CLIENT_10_ADDR2} -pki ${OPERATOR_IP}:44010 \
            -pkiCertPath /root/operator-cert.pem >> ${CLIENT_10_PATH}/log.evaluation

    else if [ "${EVAL_SYSTEM}" == "pung" ]; then

        # Run ten Pung clients.

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_01_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_01} -p ${CLIENT_01_PARTNER} -x ${CLIENT_01_PUNG_SHARED_SECRET} \
            -h ${CLIENT_01_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_01_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_02_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_02} -p ${CLIENT_02_PARTNER} -x ${CLIENT_02_PUNG_SHARED_SECRET} \
            -h ${CLIENT_02_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_02_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_03_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_03} -p ${CLIENT_03_PARTNER} -x ${CLIENT_03_PUNG_SHARED_SECRET} \
            -h ${CLIENT_03_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_03_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_04_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_04} -p ${CLIENT_04_PARTNER} -x ${CLIENT_04_PUNG_SHARED_SECRET} \
            -h ${CLIENT_04_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_04_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_05_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_05} -p ${CLIENT_05_PARTNER} -x ${CLIENT_05_PUNG_SHARED_SECRET} \
            -h ${CLIENT_05_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_05_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_06_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_06} -p ${CLIENT_06_PARTNER} -x ${CLIENT_06_PUNG_SHARED_SECRET} \
            -h ${CLIENT_06_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_06_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_07_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_07} -p ${CLIENT_07_PARTNER} -x ${CLIENT_07_PUNG_SHARED_SECRET} \
            -h ${CLIENT_07_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_07_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_08_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_08} -p ${CLIENT_08_PARTNER} -x ${CLIENT_08_PUNG_SHARED_SECRET} \
            -h ${CLIENT_08_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_08_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_09_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_09} -p ${CLIENT_09_PARTNER} -x ${CLIENT_09_PUNG_SHARED_SECRET} \
            -h ${CLIENT_09_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_09_PATH}/log.evaluation

        echo "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> ${CLIENT_10_PATH}/log.evaluation
        /root/pung-client -e /tmp/collect -n ${CLIENT_10} -p ${CLIENT_10_PARTNER} -x ${CLIENT_10_PUNG_SHARED_SECRET} \
            -h ${CLIENT_10_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> ${CLIENT_10_PATH}/log.evaluation

    else if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then

        # Run ten client components of Vuvuzela.

        echo "\n" >> ${CLIENT_01_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_01}.conf \
            -peer ${CLIENT_01_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_01_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_02_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_02}.conf \
            -peer ${CLIENT_02_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_02_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_03_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_03}.conf \
            -peer ${CLIENT_03_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_03_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_04_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_04}.conf \
            -peer ${CLIENT_04_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_04_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_05_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_05}.conf \
            -peer ${CLIENT_05_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_05_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_06_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_06}.conf \
            -peer ${CLIENT_06_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_06_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_07_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_07}.conf \
            -peer ${CLIENT_07_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_07_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_08_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_08}.conf \
            -peer ${CLIENT_08_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_08_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_09_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_09}.conf \
            -peer ${CLIENT_09_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_09_PATH}/log.evaluation

        echo "\n" >> ${CLIENT_10_PATH}/log.evaluation
        /root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe /tmp/collect -conf /root/vuvuzela-confs/${CLIENT_10}.conf \
            -peer ${CLIENT_10_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> ${CLIENT_10_PATH}/log.evaluation

    fi

fi


# Wait for metrics collector to exit.
wait


# Reset tc configuration.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc del dev ${NET_DEVICE} root
fi


# Upload result files to GCloud bucket.

if ([ "${TYPE_OF_NODE}" == "server" ] || [ "${TYPE_OF_NODE}" == "coordinator" ]); then

    /usr/bin/gsutil -m cp ${CLIENT_01_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_01}.evaluation
    /usr/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/

else if [ "${TYPE_OF_NODE}" == "client" ]; then

    /usr/bin/gsutil -m cp ${CLIENT_01_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_01}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_02_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_02}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_03_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_03}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_04_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_04}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_05_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_05}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_06_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_06}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_07_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_07}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_08_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_08}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_09_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_09}.evaluation
    /usr/bin/gsutil -m cp ${CLIENT_10_PATH}/log.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/log_${CLIENT_10}.evaluation
    /usr/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/

fi

# Mark worker as finished at operator.
curl --cacert /root/operator-cert.pem --request PUT https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/finished
