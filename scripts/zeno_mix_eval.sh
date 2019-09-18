#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system zeno -server -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run zeno as mix.
/root/zeno -eval -killMixesInRound ${KILL_ZENO_MIXES_IN_ROUND} -metricsPipe ${METRICS_PIPE} -mix -name ${CLIENT} \
    -partner ${PARTNER} -msgPublicAddr ${LISTEN_IP}:33000 -msgLisAddr ${LISTEN_IP}:33000 -pkiLisAddr ${LISTEN_IP}:44000 \
    -pki ${ZENO_PKI_IP}:33000 -pkiCertPath /root/operator-cert.pem > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
