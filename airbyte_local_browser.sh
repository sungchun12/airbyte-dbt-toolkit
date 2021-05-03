#!/bin/bash

# remove anything running on port 8000 and 8001 for lingering airbyte webserver deployment
lsof -i:8000 -i:8001 -Fp | sed 's/^p//' | xargs kill -9

# create ssh tunnel to run airbyte webserver in the local browser
gcloud beta compute ssh $INSTANCE_NAME \
--project=$GOOGLE_CLOUD_PROJECT_ID \
--zone $YOUR_GCP_ZONE \
-- -L 8000:localhost:8000 -L 8001:localhost:8001 -N -f

# open the airbyte webserver in your local browser
open http://localhost:8000/