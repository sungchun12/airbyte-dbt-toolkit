#!/bin/bash

# setup env vars
export YOUR_GCP_PROJECT="dbt-demos-sung"
export GOOGLE_CLOUD_PROJECT_ID="dbt-demos-sung"
export INSTANCE_NAME="airbyte-demo"
export YOUR_GCP_ZONE="us-central1-a"
export GOOGLE_APPLICATION_CREDENTIALS="../service_account.json"

# terratest configs
export SKIP_teardown=true
export SKIP_cleanup_image=true