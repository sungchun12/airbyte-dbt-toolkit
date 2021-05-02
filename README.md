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

## Download Gruntwork Utilities

Get a [GitHub Access Token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

- enable scopes: `public_repo`

`export GITHUB_OAUTH_TOKEN="(your secret token)"`

```bash
# install gruntwork-install
curl -LsS https://raw.githubusercontent.com/gruntwork-io/gruntwork-installer/master/bootstrap-gruntwork-installer.sh | bash /dev/stdin --version v0.0.22

# install terratest_log_parser
gruntwork-install --binary-name 'terratest_log_parser' --repo 'https://github.com/gruntwork-io/terratest' --tag 'v0.13.13'


# run tests and then parse logs for human readability
go test -timeout 30m | tee test_output.log
terratest_log_parser -testlog test_output.log -outputdir test_output

```
