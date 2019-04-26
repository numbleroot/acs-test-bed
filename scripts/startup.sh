#!/usr/bin/env bash

# Retrieve metadata required for operation.
EXTERNAL_IP=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip -H "Metadata-Flavor: Google")
TYPE_OF_NODE=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/typeOfNode -H "Metadata-Flavor: Google")
EVAL_SCRIPT_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/evalScriptToPull -H "Metadata-Flavor: Google")
BINARY_TO_PULL=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/binaryToPull -H "Metadata-Flavor: Google")
TC_CONFIG=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tcConfig -H "Metadata-Flavor: Google")

# Pull files from GCloud bucket.
/snap/bin/gsutil cp gs://acs-eval/${EVAL_SCRIPT_TO_PULL} ~/${EVAL_SCRIPT_TO_PULL}
/snap/bin/gsutil cp gs://acs-eval/${BINARY_TO_PULL} ~/${BINARY_TO_PULL}

# Make the downloaded binary executable.
chmod 0700 ~/${BINARY_TO_PULL}

# Hand over to evaluation script.
EXTERNAL_IP=${EXTERNAL_IP} TYPE_OF_NODE=${TYPE_OF_NODE} TC_CONFIG=${TC_CONFIG} /bin/bash ~/${EVAL_SCRIPT_TO_PULL}