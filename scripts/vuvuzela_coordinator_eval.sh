#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system vuvuzela -server -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run coordinator component of vuvuzela.
/root/vuvuzela-coordinator -eval -metricsPipe ${METRICS_PIPE} -addr ${LISTEN_IP}:33000 \
    -wait 10s -pki /root/vuvuzela-confs/pki.conf > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
