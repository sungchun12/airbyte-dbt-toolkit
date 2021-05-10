#!/bin/bash

# setup env vars
export GOOGLE_CLOUD_PROJECT_ID="dbt-demos-sung"
export INSTANCE_NAME="airbyte-demo"
export YOUR_GCP_ZONE="us-central1-a"
export GOOGLE_APPLICATION_CREDENTIALS="../service_account.json"

# setup terraform deployment environment variables
# https://www.terraform.io/docs/cli/config/environment-variables.html
export TF_VAR_credentials="../service_account.json"
export TF_VAR_project="dbt-demos-sung"
export TF_VAR_location="us-central1"
export TF_VAR_subnetwork="projects/$TF_VAR_project/regions/$TF_VAR_location/subnetworks/default"
export TF_VAR_zone="us-central1-a"
export TF_VAR_service_account_email="packer@dbt-demos-sung.iam.gserviceaccount.com"
export TF_VAR_version_label="demo"
export TF_VAR_image="" # default to blank as terratest will dynamically create this variable

# terratest configs for faster, local testing
# export SKIP_teardown=true
# export SKIP_build_image=true
# export SKIP_cleanup_image=true