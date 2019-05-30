#!/usr/bin/env bash

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

# Pull TLS certificates of PKI.
/snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki-${RESULT_FOLDER}.pem /root/cert.pem
chmod 0644 /root/cert.pem

tried=0
while [ ! -e /root/cert.pem ] && [ "${tried}" -lt 20 ]; do

    printf "Failed to pull cert_zeno-pki-${RESULT_FOLDER}.pem, sleeping 1 second\n"
    ls -lah /root/

    sleep 1

    # Reattempt to pull TLS certificates of PKI.
    /snap/bin/gsutil cp gs://acs-eval/cert_zeno-pki-${RESULT_FOLDER}.pem /root/cert.pem
    chmod 0644 /root/cert.pem
done

if [ "${tried}" -eq 20 ]; then
    printf "Waited 20 seconds for cert_zeno-pki-${RESULT_FOLDER}.pem to be downloaded, no success, shutting down\n"
    poweroff
fi

# Run metrics collector sidecar in background.
/root/collector -system zeno -server -pipe /tmp/collect -metricsPath /root/ &

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Add iptables rules to be able to count number of transferred
# bytes for evaluation and initialize them to zero.
iptables -A INPUT -p tcp --dport 33000
iptables -A OUTPUT -p tcp --dport 33000
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run zeno as mix.
/root/zeno -eval -killMixesInRound ${KILL_ZENO_MIXES_IN_ROUND} -metricsPipe /tmp/collect -mix -name "${NAME_OF_NODE}" -partner "${PARTNER_OF_NODE}" \
    -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr ${LISTEN_IP}:33000 -pkiLisAddr ${LISTEN_IP}:44000 \
    -pki ${PKI_IP}:33000 -pkiCertPath "/root/cert.pem" > /root/log.evaluation

# Wait for metrics collector to exit.
wait

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/*.evaluation gs://acs-eval/${RESULT_FOLDER}/servers/${NAME_OF_NODE}_${LISTEN_IP}/
