#!/bin/bash
echo "***********************"
echo "Removing anything running on port 8000 and 8001 for lingering airbyte webserver deployment"
echo "***********************"
lsof -i:8000 -i:8001 -Fp | sed 's/^p//' | xargs kill -9

echo "***********************"
echo "Creating ssh tunnel to run airbyte webserver in the local browser"
echo "***********************"
gcloud beta compute ssh $TF_VAR_name \
--project=$TF_VAR_project \
--zone $TF_VAR_zone \
-- -L 8000:localhost:8000 -L 8001:localhost:8001 -N -f

echo "***********************"
echo "Opening the airbyte webserver in your local browser"
echo "***********************"
open http://localhost:8000/

echo "***********************"
echo "You should now see your airbyte webserver in your local browser!"
echo "***********************"
