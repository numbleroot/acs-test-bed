#!/usr/bin/env bash

# Run metrics collector sidecar in background.
/root/collector -system pung -server -pipe ${METRICS_PIPE} -metricsPath ${METRICS_PATH}/ &

# TODO: Make multi-process.

# Run pung as server.
/root/pung-server -e 30 -n 1 -w 1 -i ${LISTEN_IP} -s 33000 -k 1 -t e -d 2 -b 0 -m 500 > ${METRICS_PATH}/log.evaluation

# Wait for metrics collector to exit.
wait
