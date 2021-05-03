package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Occasionally, a Packer build may fail due to intermittent issues (e.g., brief network outage or GCP issue). We try
// to make our tests resilient to that by specifying those known common errors here and telling our builds to retry if
// they hit those errors.
var DefaultRetryablePackerErrors = map[string]string{
	"Script disconnected unexpectedly":                                                 "Occasionally, Packer seems to lose connectivity to GCP, perhaps due to a brief network outage",
	"can not open /var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_xenial_InRelease": "Occasionally, apt-get fails on ubuntu to update the cache",
}
var DefaultTimeBetweenPackerRetries = 15 * time.Second

const DefaultMaxPackerRetries = 3

// An example of how to test the Terraform module in examples/terraform-ssh-example using Terratest. The test also
// shows an example of how to break a test down into "stages" so you can skip stages by setting environment variables
// (e.g., skip stage "teardown" by setting the environment variable "SKIP_teardown=true"), which speeds up iteration
// when running this test over and over again locally.
func TestTerraformAirbyteDemo(t *testing.T) {
	t.Parallel()

	workingTerraformDir := test_structure.CopyTerraformFolderToTemp(t, "../", "terraform")
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.

	// The folder where we have our packer code
	workingPackerDir := "../packer/"

	// Get the Project Id to use
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// At the end of the test, delete the Image
	defer test_structure.RunTestStage(t, "cleanup_image", func() {
		imageName := test_structure.LoadString(t, workingPackerDir, "imageName")
		deleteImage(t, projectID, imageName)
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, workingTerraformDir)
		terraform.Destroy(t, terraformOptions)
	})

	// Build the Image for the web app
	test_structure.RunTestStage(t, "build_image", func() {
		buildImage(t, projectID, workingPackerDir)
	})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: workingTerraformDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"image": test_structure.LoadString(t, workingPackerDir, "imageName"),
		},
	})

	// Deploy the terraform infrastructure
	test_structure.RunTestStage(t, "setup", func() {
		// terraformOptions, keyPair := configureTerraformOptions(t, workingTerraformDir)

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, workingTerraformDir, terraformOptions)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		terraform.InitAndApply(t, terraformOptions)
	})

	// Make sure we can SSH to the public Instance directly from the public Internet and the private Instance by using
	// the public Instance as a jump host
	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, workingTerraformDir)

		testServiceAccountRoles(t, terraformOptions)
		testComputeEngineId(t, terraformOptions)
		testBigQueryDatasetId(t, terraformOptions)
		testdbtServiceAccountEmail(t, terraformOptions, projectID)
	})

}

// Build the Packer Image
func buildImage(t *testing.T, projectID string, workingPackerDir string) {

	// Pick a random GCP zone to test in. This helps ensure your code works in all regions.
	// On October 22, 2018, GCP launched the asia-east2 region, which promptly failed all our tests, so blacklist asia-east2.
	zone := gcp.GetRandomZone(t, projectID, nil, nil, []string{"asia-east2"})
	airbyte_build_script := workingPackerDir + "airbyte_build.sh"

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: workingPackerDir + "airbyte_gce_image.pkr.hcl",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"gcp_project_id":       projectID,
			"gcp_zone":             zone,
			"airbyte_build_script": airbyte_build_script,
		},

		// Configure retries for intermittent errors
		RetryableErrors:    DefaultRetryablePackerErrors,
		TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
		MaxRetries:         DefaultMaxPackerRetries,
	}

	// Save the Packer Options so future test stages can use them
	test_structure.SavePackerOptions(t, workingPackerDir, packerOptions)

	// Make sure the Packer build completes successfully
	imageName := packer.BuildArtifact(t, packerOptions)

	// Save the imageName as a string so future test stages can use them
	test_structure.SaveString(t, workingPackerDir, "imageName", imageName)
}

// Delete the Packer Image
func deleteImage(t *testing.T, projectID string, imageName string) {
	// Load the Image ID saved by the earlier build_image stage
	image := gcp.FetchImage(t, projectID, imageName)
	image.DeleteImage(t)
}

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

func testBigQueryDatasetId(t *testing.T, terraformOptions *terraform.Options) {

	expected_airbyte_dataset_id := "projects/dbt-demos-sung/datasets/airbyte_dataset" //TODO: change to dynamic based on env vars

	output := terraform.Output(t, terraformOptions, "airbyte_dataset_id")
	assert.Equal(t, expected_airbyte_dataset_id, output)
}

func testdbtServiceAccountEmail(t *testing.T, terraformOptions *terraform.Options, projectID string) {
	iam_domain_name := ".iam.gserviceaccount.com"

	expected_service_account_email := "dbt-service-account" + "@" + projectID + iam_domain_name //TODO: change to dynamic based on env vars

	output := terraform.Output(t, terraformOptions, "service-account-dbt-email")
	assert.Equal(t, expected_service_account_email, output)
}

//TODO: test that I can NOT access my compute instance from the public internet by simply hitting the URL from my browser

//TODO: test that I can ssh into the instance
