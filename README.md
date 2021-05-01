# airbyte-dbt-toolkit

Deploy an airbyte instance with dbt models on day 1

## Setup Dev Environment

```bash
# config service accounts
gcloud auth activate-service-account --key-file service_account.json

# setup dev environment variables in interactive shell
source config.sh

```

## Test terraform deployment

```bash

# change into terraform directory
cd terraform/

# skips tearing down resources by default
# to tear down resources after testing run: `unset SKIP_teardown`

# expected terminal output:
# --- PASS: TestTerraformAirbyteDemo (53.76s)
# PASS
# ok      example.com/m/terraform 54.057s
go test -v

```
