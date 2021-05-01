#!/bin/bash

# create ssh tunnel to run airbyte webserver in the local browser
gcloud beta compute ssh $INSTANCE_NAME \
--project=$YOUR_GCP_PROJECT \
--zone $YOUR_GCP_ZONE \
-- -L 8000:localhost:8000 -L 8001:localhost:8001 -N -f