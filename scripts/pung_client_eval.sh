#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -client -pipe ${METRICS_PIPE} -metricsPath ${METRICS_PATH}/ &

# Run pung as client.
/root/pung-client -e ${METRICS_PIPE} -n ${CLIENT} -p ${PARTNER} -x ${SHARED_SECRET} \
    -h ${SERVER_IP}:33000 -r 30 -k 1 -s 1 -t e -d 2 -b 0 > ${METRICS_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
