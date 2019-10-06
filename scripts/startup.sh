#!/usr/bin/env bash

sleep 1

# Make sure the application ports we are going to
# use for any component of the ACS we are about to
# evaluate are blocked off from "randomly binding"
# applications (=> part of the reserved pool).
sysctl -w net.ipv4.ip_local_reserved_ports=33001-33010,44001-44010

sleep 15

# Heavily increase limit on open file descriptors and
# connections per socket in order to be able to keep
# lots of connections open.
sysctl -w fs.file-max=1048575
sysctl -w net.core.somaxconn=8192
ulimit -n 1048575


# Retrieve metadata required for operation.

OPERATOR_IP=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/operatorIP" -H "Metadata-Flavor: Google")
EXP_ID=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/expID" -H "Metadata-Flavor: Google")
NAME_OF_NODE=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/nameOfNode" -H "Metadata-Flavor: Google")
EVAL_SYSTEM=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalSystem" -H "Metadata-Flavor: Google")
NUM_CLIENTS=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/numClients" -H "Metadata-Flavor: Google")
RESULT_FOLDER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/resultFolder" -H "Metadata-Flavor: Google")

LISTEN_IP=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip" -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode" -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull" -H "Metadata-Flavor: Google")

PUNG_SERVER_IP=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/pungServerIP" -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig" -H "Metadata-Flavor: Google")
KILL_ZENO_MIXES_IN_ROUND=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/killZenoMixesInRound" -H "Metadata-Flavor: Google")
PUNG_CLIENTS_PER_PROC=$(( NUM_CLIENTS / 10))


# Prepare to evaluate up to ten clients in case
# this is a clients machine. In case this is a
# server machine, values for CLIENT_01 will be
# exclusively used.

CLIENT_01=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client01" -H "Metadata-Flavor: Google")
CLIENT_01_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner01" -H "Metadata-Flavor: Google")
CLIENT_01_ADDR1="${LISTEN_IP}:33001"
CLIENT_01_ADDR2="${LISTEN_IP}:44001"
CLIENT_01_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33001"
CLIENT_01_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}12"

CLIENT_02=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client02" -H "Metadata-Flavor: Google")
CLIENT_02_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner02" -H "Metadata-Flavor: Google")
CLIENT_02_ADDR1="${LISTEN_IP}:33002"
CLIENT_02_ADDR2="${LISTEN_IP}:44002"
CLIENT_02_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33002"
CLIENT_02_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}12"

CLIENT_03=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client03" -H "Metadata-Flavor: Google")
CLIENT_03_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner03" -H "Metadata-Flavor: Google")
CLIENT_03_ADDR1="${LISTEN_IP}:33003"
CLIENT_03_ADDR2="${LISTEN_IP}:44003"
CLIENT_03_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33003"
CLIENT_03_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}34"

CLIENT_04=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client04" -H "Metadata-Flavor: Google")
CLIENT_04_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner04" -H "Metadata-Flavor: Google")
CLIENT_04_ADDR1="${LISTEN_IP}:33004"
CLIENT_04_ADDR2="${LISTEN_IP}:44004"
CLIENT_04_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33004"
CLIENT_04_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}34"

CLIENT_05=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client05" -H "Metadata-Flavor: Google")
CLIENT_05_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner05" -H "Metadata-Flavor: Google")
CLIENT_05_ADDR1="${LISTEN_IP}:33005"
CLIENT_05_ADDR2="${LISTEN_IP}:44005"
CLIENT_05_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33005"
CLIENT_05_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}56"

CLIENT_06=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client06" -H "Metadata-Flavor: Google")
CLIENT_06_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner06" -H "Metadata-Flavor: Google")
CLIENT_06_ADDR1="${LISTEN_IP}:33006"
CLIENT_06_ADDR2="${LISTEN_IP}:44006"
CLIENT_06_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33006"
CLIENT_06_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}56"

CLIENT_07=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client07" -H "Metadata-Flavor: Google")
CLIENT_07_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner07" -H "Metadata-Flavor: Google")
CLIENT_07_ADDR1="${LISTEN_IP}:33007"
CLIENT_07_ADDR2="${LISTEN_IP}:44007"
CLIENT_07_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33007"
CLIENT_07_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}78"

CLIENT_08=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client08" -H "Metadata-Flavor: Google")
CLIENT_08_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner08" -H "Metadata-Flavor: Google")
CLIENT_08_ADDR1="${LISTEN_IP}:33008"
CLIENT_08_ADDR2="${LISTEN_IP}:44008"
CLIENT_08_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33008"
CLIENT_08_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}78"

CLIENT_09=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client09" -H "Metadata-Flavor: Google")
CLIENT_09_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner09" -H "Metadata-Flavor: Google")
CLIENT_09_ADDR1="${LISTEN_IP}:33009"
CLIENT_09_ADDR2="${LISTEN_IP}:44009"
CLIENT_09_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33009"
CLIENT_09_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}910"

CLIENT_10=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/client10" -H "Metadata-Flavor: Google")
CLIENT_10_PARTNER=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/partner10" -H "Metadata-Flavor: Google")
CLIENT_10_ADDR1="${LISTEN_IP}:33010"
CLIENT_10_ADDR2="${LISTEN_IP}:44010"
CLIENT_10_PUNG_SERVER_ADDR="${PUNG_SERVER_IP}:33010"
CLIENT_10_PUNG_SHARED_SECRET="${NAME_OF_NODE}${LISTEN_IP}910"


# Prepare FIFO pipe for system and collector IPC.
mkfifo /tmp/collect01 /tmp/collect02 /tmp/collect03 /tmp/collect04 /tmp/collect05 \
    /tmp/collect06 /tmp/collect07 /tmp/collect08 /tmp/collect09 /tmp/collect10
chmod 0600 /tmp/collect01 /tmp/collect02 /tmp/collect03 /tmp/collect04 /tmp/collect05 \
    /tmp/collect06 /tmp/collect07 /tmp/collect08 /tmp/collect09 /tmp/collect10


# Register with operator for current experiment.
printf "Will call /experiments/${EXP_ID}/workers/${NAME_OF_NODE}/register as ${NAME_OF_NODE}@${LISTEN_IP}.\n"
curl --cacert /root/operator-cert.pem --request PUT --header "content-type: application/json" --data-binary "{
    \"address\": \"${CLIENT_01_ADDR1}\"
}" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/register


# Pull files from GCloud bucket.
/usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
/usr/bin/gsutil cp gs://acs-eval/collector /root/collector

tried=0
while ([ ! -e /root/${BINARY_TO_PULL} ] || [ ! -e /root/collector ]) && [ "${tried}" -lt 20 ]; do

    printf "Failed to pull required files from GCloud bucket, sleeping 1 second...\n"
    ls -lah /root/

    sleep 1

    # Reattempt to pull files from GCloud bucket.
    /usr/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} /root/${BINARY_TO_PULL}
    /usr/bin/gsutil cp gs://acs-eval/collector /root/collector

    tried=$(( tried + 1 ))
done

if [ "${tried}" -eq 20 ]; then

    printf "Waited 20 seconds for required experiment files to be downloaded, no success, shutting down.\n"

    # Inform operator about failure to initialize.
    printf "Will call /experiments/${EXP_ID}/workers/${NAME_OF_NODE}/failed as ${NAME_OF_NODE}@${LISTEN_IP}.\n"
    curl --cacert /root/operator-cert.pem --request PUT --header "content-type: application/json" --data-binary "{
        \"failure\": \"waited 20 seconds for required experiment files to be downloaded from Storage bucket, no success, shutting down\"
    }" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/failed

    poweroff
fi

if [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then
    
    sleep 10

    printf "This is a Vuvuzela experiment, pull 'gs://acs-eval/vuvuzela-confs/pki.conf' as well.\n"
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
printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_01}_log.evaluation
printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_01}_log.evaluation
printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_01}_log.evaluation
printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_01}_log.evaluation
printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_01}_log.evaluation
printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_01}_log.evaluation
printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_01}_log.evaluation

if [ "${TYPE_OF_NODE}" == "client" ]; then

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_02}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_02}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_02}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_02}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_02}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_02}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_02}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_03}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_03}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_03}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_03}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_03}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_03}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_03}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_04}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_04}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_04}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_04}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_04}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_04}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_04}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_05}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_05}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_05}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_05}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_05}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_05}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_05}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_06}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_06}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_06}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_06}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_06}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_06}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_06}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_07}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_07}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_07}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_07}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_07}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_07}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_07}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_08}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_08}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_08}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_08}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_08}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_08}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_08}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_09}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_09}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_09}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_09}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_09}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_09}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_09}_log.evaluation

    printf "Evaluating a '${TYPE_OF_NODE}' for system '${EVAL_SYSTEM}' as part of experiment '${EXP_ID}' on machine '${NAME_OF_NODE}'.\n" > /root/${CLIENT_10}_log.evaluation
    printf "Result folder: '${RESULT_FOLDER}.'\n" >> /root/${CLIENT_10}_log.evaluation
    printf "${NUM_CLIENTS} clients will participate, TC parameters set to: '%q'.\n" "${TC_CONFIG}" >> /root/${CLIENT_10}_log.evaluation
    printf "System info: '$(uname -a)'.\n" >> /root/${CLIENT_10}_log.evaluation
    printf "CPU: '$(grep ^cpu\\scores /proc/cpuinfo | uniq | awk '{print $4}') cores ($(grep -c ^processor /proc/cpuinfo) threads) as part of $(lscpu --json | grep "Model name" | awk -F \" '{print $8}')'.\n" >> /root/${CLIENT_10}_log.evaluation
    printf "Memory: '$(lsmem | grep "Total online memory" | awk '{print $4}')'.\n" >> /root/${CLIENT_10}_log.evaluation
    printf "Storage: '$(lsblk -o TYPE,SIZE | grep disk | awk '{print $2}') $(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/disks/0/type" -H "Metadata-Flavor: Google")'.\n" >> /root/${CLIENT_10}_log.evaluation

fi


sleep 5

# Signal readiness of process to experiment script.
printf "Will call /experiments/${EXP_ID}/workers/${NAME_OF_NODE}/ready as ${NAME_OF_NODE}@${LISTEN_IP}.\n"
curl --cacert /root/operator-cert.pem --request PUT --header "content-type: application/json" \
    https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/ready


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


# We will use this array to keep track of
# the status of each spawned subprocess.
declare -a PROCESS_IDS


# Run metrics collector sidecar in background.
/root/collector -system ${EVAL_SYSTEM} -typeOfNode ${TYPE_OF_NODE} -metricsPath /root/ \
    -client01 ${CLIENT_01} -pipe01 /tmp/collect01 -client02 ${CLIENT_02} -pipe02 /tmp/collect02 \
    -client03 ${CLIENT_03} -pipe03 /tmp/collect03 -client04 ${CLIENT_04} -pipe04 /tmp/collect04 \
    -client05 ${CLIENT_05} -pipe05 /tmp/collect05 -client06 ${CLIENT_06} -pipe06 /tmp/collect06 \
    -client07 ${CLIENT_07} -pipe07 /tmp/collect07 -client08 ${CLIENT_08} -pipe08 /tmp/collect08 \
    -client09 ${CLIENT_09} -pipe09 /tmp/collect09 -client10 ${CLIENT_10} -pipe10 /tmp/collect10 &
PROCESS_IDS+=($!)


if [ "${TYPE_OF_NODE}" == "server" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_01}_log.evaluation

        # Run zeno as mix.
        /root/zeno -eval -killMixesInRound ${KILL_ZENO_MIXES_IN_ROUND} -metricsPipe /tmp/collect01 -mix -name ${CLIENT_01} \
            -partner ${CLIENT_01_PARTNER} -msgPublicAddr ${CLIENT_01_ADDR1} -msgLisAddr ${CLIENT_01_ADDR1} -pkiLisAddr ${CLIENT_01_ADDR2} \
            -pki ${OPERATOR_IP}:44001 -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_01}_log.evaluation

    elif [ "${EVAL_SYSTEM}" == "pung" ]; then

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_01}_log.evaluation

        # Run Pung's server.
        /root/pung-server -e 30 -i ${LISTEN_IP} -s 33001 -n 1 -w 10 -p 0 -k 1 -t e -d 2 -b 0 -m ${PUNG_CLIENTS_PER_PROC} >> /root/${CLIENT_01}_log.evaluation

        # Force collector exit when Pung's server
        # finished its operation.
        echo "done\n" > /tmp/collect01

    elif [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then

        printf "\n" >> /root/${CLIENT_01}_log.evaluation

        # Run mix component of Vuvuzela.
        /root/vuvuzela-mix -metricsPipe /tmp/collect01 -addr ${CLIENT_01_ADDR1} -conf /root/vuvuzela-confs/${CLIENT_01}.conf \
            -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_01}_log.evaluation

    fi

elif [ "${TYPE_OF_NODE}" == "coordinator" ]; then

    printf "\n" >> /root/${CLIENT_01}_log.evaluation

    # Run coordinator component of Vuvuzela.
    /root/vuvuzela-coordinator -metricsPipe /tmp/collect01 -addr ${CLIENT_01_ADDR1} \
        -wait 10s -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_01}_log.evaluation

elif [ "${TYPE_OF_NODE}" == "client" ]; then

    if [ "${EVAL_SYSTEM}" == "zeno" ]; then

        # Run ten zeno clients.

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_01}_log.evaluation        
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect01 -client -name ${CLIENT_01} -partner ${CLIENT_01_PARTNER} \
            -msgPublicAddr ${CLIENT_01_ADDR1} -msgLisAddr ${CLIENT_01_ADDR1} -pkiLisAddr ${CLIENT_01_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_01}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_02}_log.evaluation 
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect02 -client -name ${CLIENT_02} -partner ${CLIENT_02_PARTNER} \
            -msgPublicAddr ${CLIENT_02_ADDR1} -msgLisAddr ${CLIENT_02_ADDR1} -pkiLisAddr ${CLIENT_02_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_02}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_03}_log.evaluation 
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect03 -client -name ${CLIENT_03} -partner ${CLIENT_03_PARTNER} \
            -msgPublicAddr ${CLIENT_03_ADDR1} -msgLisAddr ${CLIENT_03_ADDR1} -pkiLisAddr ${CLIENT_03_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_03}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_04}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect04 -client -name ${CLIENT_04} -partner ${CLIENT_04_PARTNER} \
            -msgPublicAddr ${CLIENT_04_ADDR1} -msgLisAddr ${CLIENT_04_ADDR1} -pkiLisAddr ${CLIENT_04_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_04}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_05}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect05 -client -name ${CLIENT_05} -partner ${CLIENT_05_PARTNER} \
            -msgPublicAddr ${CLIENT_05_ADDR1} -msgLisAddr ${CLIENT_05_ADDR1} -pkiLisAddr ${CLIENT_05_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_05}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_06}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect06 -client -name ${CLIENT_06} -partner ${CLIENT_06_PARTNER} \
            -msgPublicAddr ${CLIENT_06_ADDR1} -msgLisAddr ${CLIENT_06_ADDR1} -pkiLisAddr ${CLIENT_06_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_06}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_07}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect07 -client -name ${CLIENT_07} -partner ${CLIENT_07_PARTNER} \
            -msgPublicAddr ${CLIENT_07_ADDR1} -msgLisAddr ${CLIENT_07_ADDR1} -pkiLisAddr ${CLIENT_07_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_07}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_08}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect08 -client -name ${CLIENT_08} -partner ${CLIENT_08_PARTNER} \
            -msgPublicAddr ${CLIENT_08_ADDR1} -msgLisAddr ${CLIENT_08_ADDR1} -pkiLisAddr ${CLIENT_08_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_08}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_09}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect09 -client -name ${CLIENT_09} -partner ${CLIENT_09_PARTNER} \
            -msgPublicAddr ${CLIENT_09_ADDR1} -msgLisAddr ${CLIENT_09_ADDR1} -pkiLisAddr ${CLIENT_09_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_09}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Some zeno mixes will be terminated in round: '${KILL_ZENO_MIXES_IN_ROUND}'.\n\n" >> /root/${CLIENT_10}_log.evaluation
        /root/zeno -eval -numMsgToRecv 25 -metricsPipe /tmp/collect10 -client -name ${CLIENT_10} -partner ${CLIENT_10_PARTNER} \
            -msgPublicAddr ${CLIENT_10_ADDR1} -msgLisAddr ${CLIENT_10_ADDR1} -pkiLisAddr ${CLIENT_10_ADDR2} -pki ${OPERATOR_IP}:44001 \
            -pkiCertPath /root/operator-cert.pem >> /root/${CLIENT_10}_log.evaluation &
        PROCESS_IDS+=($!)

    elif [ "${EVAL_SYSTEM}" == "pung" ]; then

        # Run ten Pung clients.

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_01}_log.evaluation
        /root/pung-client -e /tmp/collect01 -n ${CLIENT_01} -p ${CLIENT_01_PARTNER} -x ${CLIENT_01_PUNG_SHARED_SECRET} \
            -h ${CLIENT_01_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_01}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_02}_log.evaluation
        /root/pung-client -e /tmp/collect02 -n ${CLIENT_02} -p ${CLIENT_02_PARTNER} -x ${CLIENT_02_PUNG_SHARED_SECRET} \
            -h ${CLIENT_02_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_02}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_03}_log.evaluation
        /root/pung-client -e /tmp/collect03 -n ${CLIENT_03} -p ${CLIENT_03_PARTNER} -x ${CLIENT_03_PUNG_SHARED_SECRET} \
            -h ${CLIENT_03_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_03}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_04}_log.evaluation
        /root/pung-client -e /tmp/collect04 -n ${CLIENT_04} -p ${CLIENT_04_PARTNER} -x ${CLIENT_04_PUNG_SHARED_SECRET} \
            -h ${CLIENT_04_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_04}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_05}_log.evaluation
        /root/pung-client -e /tmp/collect05 -n ${CLIENT_05} -p ${CLIENT_05_PARTNER} -x ${CLIENT_05_PUNG_SHARED_SECRET} \
            -h ${CLIENT_05_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_05}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_06}_log.evaluation
        /root/pung-client -e /tmp/collect06 -n ${CLIENT_06} -p ${CLIENT_06_PARTNER} -x ${CLIENT_06_PUNG_SHARED_SECRET} \
            -h ${CLIENT_06_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_06}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_07}_log.evaluation
        /root/pung-client -e /tmp/collect07 -n ${CLIENT_07} -p ${CLIENT_07_PARTNER} -x ${CLIENT_07_PUNG_SHARED_SECRET} \
            -h ${CLIENT_07_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_07}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_08}_log.evaluation
        /root/pung-client -e /tmp/collect08 -n ${CLIENT_08} -p ${CLIENT_08_PARTNER} -x ${CLIENT_08_PUNG_SHARED_SECRET} \
            -h ${CLIENT_08_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_08}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_09}_log.evaluation
        /root/pung-client -e /tmp/collect09 -n ${CLIENT_09} -p ${CLIENT_09_PARTNER} -x ${CLIENT_09_PUNG_SHARED_SECRET} \
            -h ${CLIENT_09_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_09}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "Pung server at: '${PUNG_SERVER_IP}', expecting ${PUNG_CLIENTS_PER_PROC} clients per process.\n\n" >> /root/${CLIENT_10}_log.evaluation
        /root/pung-client -e /tmp/collect10 -n ${CLIENT_10} -p ${CLIENT_10_PARTNER} -x ${CLIENT_10_PUNG_SHARED_SECRET} \
            -h ${CLIENT_10_PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 >> /root/${CLIENT_10}_log.evaluation &
        PROCESS_IDS+=($!)

    elif [ "${EVAL_SYSTEM}" == "vuvuzela" ]; then

        # Run ten client components of Vuvuzela.

        printf "\n" >> /root/${CLIENT_01}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect01 -conf /root/vuvuzela-confs/${CLIENT_01}.conf \
            -peer ${CLIENT_01_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_01}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_02}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect02 -conf /root/vuvuzela-confs/${CLIENT_02}.conf \
            -peer ${CLIENT_02_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_02}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_03}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect03 -conf /root/vuvuzela-confs/${CLIENT_03}.conf \
            -peer ${CLIENT_03_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_03}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_04}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect04 -conf /root/vuvuzela-confs/${CLIENT_04}.conf \
            -peer ${CLIENT_04_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_04}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_05}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect05 -conf /root/vuvuzela-confs/${CLIENT_05}.conf \
            -peer ${CLIENT_05_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_05}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_06}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect06 -conf /root/vuvuzela-confs/${CLIENT_06}.conf \
            -peer ${CLIENT_06_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_06}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_07}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect07 -conf /root/vuvuzela-confs/${CLIENT_07}.conf \
            -peer ${CLIENT_07_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_07}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_08}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect08 -conf /root/vuvuzela-confs/${CLIENT_08}.conf \
            -peer ${CLIENT_08_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_08}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_09}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect09 -conf /root/vuvuzela-confs/${CLIENT_09}.conf \
            -peer ${CLIENT_09_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_09}_log.evaluation &
        PROCESS_IDS+=($!)

        printf "\n" >> /root/${CLIENT_10}_log.evaluation
        /root/vuvuzela-client -numMsgToRecv 30 -metricsPipe /tmp/collect10 -conf /root/vuvuzela-confs/${CLIENT_10}.conf \
            -peer ${CLIENT_10_PARTNER} -pki /root/vuvuzela-confs/pki.conf >> /root/${CLIENT_10}_log.evaluation &
        PROCESS_IDS+=($!)

    fi

fi


for PROCESS_ID in "${PROCESS_IDS[@]}"; do

    # Wait for the next process to finish.
    wait -n

    # Extract its return value.
    RET_VALUE=$?

    if (( ${RET_VALUE} > 0 )); then

        # Reset tc configuration.
        if [ "${TC_CONFIG}" != "none" ]; then
            tc qdisc del dev ${NET_DEVICE} root
        fi

        # If the process returned an error code
        # tell the operator about it.
        printf "Will call /experiments/${EXP_ID}/workers/${NAME_OF_NODE}/failed as ${NAME_OF_NODE}@${LISTEN_IP}.\n"
        curl --cacert /root/operator-cert.pem --request PUT --header "content-type: application/json" --data-binary "{
            \"failure\": \"one or more client processes exited with an error code\"
        }" https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/failed
    fi

done


# Reset tc configuration.
if [ "${TC_CONFIG}" != "none" ]; then
    tc qdisc del dev ${NET_DEVICE} root
fi


# Upload result files to GCloud bucket.

if ([ "${TYPE_OF_NODE}" == "server" ] || [ "${TYPE_OF_NODE}" == "coordinator" ]); then

    printf "Uploading results as '${TYPE_OF_NODE}' to '${RESULT_FOLDER}' under '/servers/${NAME_OF_NODE}_${LISTEN_IP}'\n"

    /usr/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/

elif [ "${TYPE_OF_NODE}" == "client" ]; then

    printf "Uploading results as '${TYPE_OF_NODE}' to '${RESULT_FOLDER}' under '/clients/${NAME_OF_NODE}_${LISTEN_IP}'\n"

    /usr/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/clients/${NAME_OF_NODE}_${LISTEN_IP}/

fi

# Mark client as finished at operator.
printf "Will call /experiments/${EXP_ID}/workers/${NAME_OF_NODE}/finished as ${NAME_OF_NODE}@${LISTEN_IP}.\n"
curl --cacert /root/operator-cert.pem --request PUT --header "content-type: application/json" \
    https://${OPERATOR_IP}/internal/experiments/${EXP_ID}/workers/${NAME_OF_NODE}/finished
