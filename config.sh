#!/bin/bash

# setup env vars
export GOOGLE_CLOUD_PROJECT_ID="dbt-demos-sung"
export INSTANCE_NAME="airbyte-demo"
export YOUR_GCP_ZONE="us-central1-a"
export GOOGLE_APPLICATION_CREDENTIALS="../service_account.json"

# terratest configs for faster, local testing
export SKIP_teardown=true
# export SKIP_build_image=true
export SKIP_cleanup_image=true