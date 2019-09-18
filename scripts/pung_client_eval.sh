#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -client -pipe ${METRICS_PIPE} -metricsPath ${CLIENT_PATH}/ &

# Run pung as client.
/root/pung-client -e ${METRICS_PIPE} -n ${CLIENT} -p ${PARTNER} -x ${PUNG_SHARED_SECRET} \
    -h ${PUNG_SERVER_ADDR} -r 30 -k 1 -s 1 -t e -d 2 -b 0 > ${CLIENT_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
