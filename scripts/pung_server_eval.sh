#!/usr/bin/env bash

cat <<EOT >> /root/timely_worker_hosts.txt
127.0.0.1:1234
127.0.0.1:1235
127.0.0.1:1236
EOT

for i in $(seq -f "%05g" 1 3)
do
    # Prepare result folders.
    mkdir /root/server-${i}

    # Create all named pipes.
    rm -f /tmp/collect
    mkfifo /tmp/collect-${i}
    chmod 0600 /tmp/collect-${i}

    # Run metrics collector sidecar in background.
    /root/collector -mix -pipe /tmp/collect-${i} -metricsPath /root/server-${i}/ &
done

# Signal readiness of process to experiment script.
curl -X PUT --data "ThisNodeIsReady" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/acs-eval/initStatus -H "Metadata-Flavor: Google"

# Configure tc according to environment variable.

# Reset bytes counters right before starting pung.
iptables -Z -t filter -L INPUT
iptables -Z -t filter -L OUTPUT

# Run multiple pung servers.
/root/pung -i ${LISTEN_IP} -s 33000 -d 2 -b 0 -t b -m 200 -h /root/timely_worker_hosts.txt -n 3 -p 0 -k 64 -o h2 > /root/server-00001/log.evaluation &
/root/pung -i ${LISTEN_IP} -s 33000 -d 2 -b 0 -t b -m 200 -h /root/timely_worker_hosts.txt -n 3 -p 0 -k 64 -o h2 > /root/server-00002/log.evaluation &
/root/pung -i ${LISTEN_IP} -s 33000 -d 2 -b 0 -t b -m 200 -h /root/timely_worker_hosts.txt -n 3 -p 0 -k 64 -o h2 > /root/server-00003/log.evaluation

# Wait for all spawned processes to exit.
wait

# Reset tc configuration.

# Upload result files to GCloud bucket.
/snap/bin/gsutil -m cp /root/server-* gs://acs-eval/${RESULT_FOLDER}/servers/
