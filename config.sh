#!/bin/bash

# setup default Google Cloud credentials to run commands tied to your project
export GOOGLE_APPLICATION_CREDENTIALS="../service_account.json"

# setup terraform deployment environment variables
# these variables will work consistently across terraform and terratest deployments
# https://www.terraform.io/docs/cli/config/environment-variables.html
export TF_VAR_credentials="../service_account.json"
export TF_VAR_project="dbt-demos-sung"
export TF_VAR_location="us-central1"
export TF_VAR_subnetwork="projects/$TF_VAR_project/regions/$TF_VAR_location/subnetworks/default"
export TF_VAR_zone="us-central1-a"
export TF_VAR_service_account_email="packer@dbt-demos-sung.iam.gserviceaccount.com"
export TF_VAR_version_label="demo"
export TF_VAR_image="" # default to blank as terratest will dynamically create this variable
export TF_VAR_name="airbyte-demo-sung"
export TF_VAR_airbyte_dataset_id="airbyte_dataset"

# setup packer variables to inherit terraform environment variables if applicable
export PKR_VAR_project=$TF_VAR_project
export PKR_VAR_zone=$TF_VAR_zone
export PKR_VAR_airbyte_build_script="airbyte_build.sh"
export PKR_VAR_airbyte_source_image="debian-10-buster-v20210420"
export PKR_VAR_airbyte_ssh_username="packer"

# terratest configs for faster, local testing
export SKIP_teardown=true
export SKIP_build_image=true
export SKIP_cleanup_image=true

# unset SKIP env variables
# unset SKIP_teardown SKIP_build_image SKIP_cleanup_image
