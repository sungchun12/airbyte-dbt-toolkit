package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-ssh-example using Terratest. The test also
// shows an example of how to break a test down into "stages" so you can skip stages by setting environment variables
// (e.g., skip stage "teardown" by setting the environment variable "SKIP_teardown=true"), which speeds up iteration
// when running this test over and over again locally.
func TestTerraformAirbyteDemo(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "terraform")
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: exampleFolder,
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.Destroy(t, terraformOptions)
	})

	// Deploy the example
	test_structure.RunTestStage(t, "setup", func() {
		// terraformOptions, keyPair := configureTerraformOptions(t, exampleFolder)

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		terraform.InitAndApply(t, terraformOptions)
	})

	// Make sure we can SSH to the public Instance directly from the public Internet and the private Instance by using
	// the public Instance as a jump host
	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)

		testServiceAccountRoles(t, terraformOptions)
		testComputeEngineId(t, terraformOptions)
	})

}

//TODO: test if the service account created has at least the two bigquery roles in scope
//? does this function ahve to be lowercase for it to be called by another function?
func testServiceAccountRoles(t *testing.T, terraformOptions *terraform.Options) {

	expected_permissions_list := string("[roles/bigquery.dataEditor roles/bigquery.user]")

	output := terraform.Output(t, terraformOptions, "dbt-iam-permissions-list")
	assert.Equal(t, expected_permissions_list, output) //TODO: change to contain assert
}

func testComputeEngineId(t *testing.T, terraformOptions *terraform.Options) {

	expected_compute_engine_id := "projects/dbt-demos-sung/zones/us-central1-a/instances/airbyte-demo" //TODO: change to dynamic based on env vars

	output := terraform.Output(t, terraformOptions, "compute_engine_id")
	assert.Equal(t, expected_compute_engine_id, output)
}

//? Should I have different capitalized test functions. Need to see other examples

//TODO: test if the service account terraform output matches the expected email

//TODO: test if the bigquery dataset id matches an expected id

//TODO: test that I can NOT access my compute instance from the public internet by simply hitting the URL from my browser

//TODO: test that I can ssh into the instance
