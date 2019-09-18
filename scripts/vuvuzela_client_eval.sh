#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system vuvuzela -client -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run client component of vuvuzela.
/root/vuvuzela-client -eval -numMsgToRecv 25 -metricsPipe ${METRICS_PIPE} \
    -conf /root/vuvuzela-confs/${CLIENT}.conf -peer ${PARTNER} \
    -pki /root/vuvuzela-confs/pki.conf > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
