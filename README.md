# airbyte-dbt-toolkit

Deploy an airbyte instance with dbt models on day 1

## Setup Dev Environment

```bash
# config service accounts
gcloud auth activate-service-account --key-file service_account.json

# setup dev environment variables in interactive shell
source config.sh

```

## Deploy and test terraform and packer resources

> Estimated time to complete: 5 minutes

TODO: add a table of resources deployed and why

```bash

# change into test directory
cd test/

# skips tearing down resources by default
# to tear down resources after testing run: `unset SKIP_teardown SKIP_cleanup_image`
# expected terminal output:
# --- PASS: TestTerraformAirbyteDemo (53.76s)
# PASS
# ok      example.com/m/terraform 54.057s
go test -v

# run airbyte webserver in your local browser through an ssh tunnel into the airbyte virtual machine
source ../airbyte_local_browser.sh

open http://localhost:8000/

```
