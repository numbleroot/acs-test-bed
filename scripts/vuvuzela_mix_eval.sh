#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system vuvuzela -server -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run mix component of vuvuzela.
/root/vuvuzela-mix -eval -metricsPipe ${METRICS_PIPE} -conf /root/vuvuzela-confs/${CLIENT}.conf \
    -pki /root/vuvuzela-confs/pki.conf > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
