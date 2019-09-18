#!/usr/bin/env bash

CLIENT
PARTNER
LISTEN_IP
ADDR1
ADDR2
CLIENT_PATH
METRICS_PIPE
ZENO_PKI_IP
PUNG_SERVER_ADDR
PUNG_CLIENTS_PER_PROC
KILL_ZENO_MIXES_IN_ROUND
PUNG_SHARED_SECRET


# Run metrics collector sidecar in background.
/root/collector -system zeno -client -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run zeno as client.
/root/zeno -eval -numMsgToRecv 25 -metricsPipe ${METRICS_PIPE} -client -name ${CLIENT} -partner ${PARTNER} \
    -msgPublicAddr ${ADDR1} -msgLisAddr ${ADDR1} -pkiLisAddr ${ADDR2} -pki ${ZENO_PKI_IP}:33000 \
    -pkiCertPath /root/operator-cert.pem > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
