package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
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

	// Get all the terraform variables to pass through
	TF_VAR_credentials := os.Getenv("TF_VAR_credentials")
	TF_VAR_project := os.Getenv("TF_VAR_project")
	TF_VAR_location := os.Getenv("TF_VAR_location")
	TF_VAR_subnetwork := os.Getenv("TF_VAR_subnetwork")
	TF_VAR_zone := os.Getenv("TF_VAR_zone")
	TF_VAR_service_account_email := os.Getenv("TF_VAR_service_account_email")
	TF_VAR_version_label := os.Getenv("TF_VAR_version_label")
	TF_VAR_name := os.Getenv("TF_VAR_name")
	TF_VAR_airbyte_dataset_id := os.Getenv("TF_VAR_airbyte_dataset_id")

	// At the end of the test, delete the Image
	defer test_structure.RunTestStage(t, "cleanup_image", func() {
		imageName := test_structure.LoadString(t, workingPackerDir, "imageName")
		deleteImage(t, TF_VAR_project, imageName)
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, workingTerraformDir)
		terraform.Destroy(t, terraformOptions)
	})

	// Build the Image for the web app
	test_structure.RunTestStage(t, "build_image", func() {
		buildImage(t, TF_VAR_project, workingPackerDir)
		logger.Log(t, "--- PASS: build_image")
	})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: workingTerraformDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"credentials":           TF_VAR_credentials,
			"project":               TF_VAR_project,
			"location":              TF_VAR_location,
			"subnetwork":            TF_VAR_subnetwork,
			"zone":                  TF_VAR_zone,
			"service_account_email": TF_VAR_service_account_email,
			"version_label":         TF_VAR_version_label,
			"name":                  TF_VAR_name,
			"airbyte_dataset_id":    TF_VAR_airbyte_dataset_id,
			"image":                 test_structure.LoadString(t, workingPackerDir, "imageName"),
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
		testdbtServiceAccountEmail(t, terraformOptions)
		testSSHToPublicHost(t, terraformOptions)
	})

}

// Build the Packer Image
func buildImage(t *testing.T, projectID string, workingPackerDir string) {
	// Variables to pass to our Packer build using -var options
	PKR_VAR_project := os.Getenv("PKR_VAR_project")
	PKR_VAR_zone := os.Getenv("PKR_VAR_zone")
	PKR_VAR_airbyte_build_script := os.Getenv("PKR_VAR_airbyte_build_script")
	airbyte_build_script := workingPackerDir + PKR_VAR_airbyte_build_script

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: workingPackerDir + "airbyte_gce_image.pkr.hcl",

		Vars: map[string]string{
			"project":              PKR_VAR_project,
			"zone":                 PKR_VAR_zone,
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
	logger.Log(t, "--- PASS: testServiceAccountRoles")
}

func testComputeEngineId(t *testing.T, terraformOptions *terraform.Options) {
	// get all the vars from the terraformOptions
	project := terraformOptions.Vars["project"].(string)
	zone := terraformOptions.Vars["zone"].(string)
	name := terraformOptions.Vars["name"].(string)
	expected_compute_engine_id := "projects/" + project + "/zones/" + zone + "/instances/" + name

	output := terraform.Output(t, terraformOptions, "compute_engine_id")
	assert.Equal(t, expected_compute_engine_id, output)
	logger.Log(t, "--- PASS: testComputeEngineId") // TODO: make this a conditional log
}

func testBigQueryDatasetId(t *testing.T, terraformOptions *terraform.Options) {
	// get all the vars from the terraformOptions
	project := terraformOptions.Vars["project"].(string)
	airbyte_dataset_id := terraformOptions.Vars["airbyte_dataset_id"].(string)
	expected_airbyte_dataset_id := "projects/" + project + "/datasets/" + airbyte_dataset_id

	output := terraform.Output(t, terraformOptions, "airbyte_dataset_id")
	assert.Equal(t, expected_airbyte_dataset_id, output)
	logger.Log(t, "--- PASS: testBigQueryDatasetId")
}

func testdbtServiceAccountEmail(t *testing.T, terraformOptions *terraform.Options) {
	iam_domain_name := ".iam.gserviceaccount.com"
	project := terraformOptions.Vars["project"].(string)
	expected_service_account_email := "dbt-service-account" + "@" + project + iam_domain_name

	output := terraform.Output(t, terraformOptions, "service-account-dbt-email")
	assert.Equal(t, expected_service_account_email, output)
	logger.Log(t, "--- PASS: testdbtServiceAccountEmail")
}

// test that I can ssh into the airbyte demo instance
func testSSHToPublicHost(t *testing.T, terraformOptions *terraform.Options) {
	//get the GCP instance
	project := terraformOptions.Vars["project"].(string)
	instanceName := terraformOptions.Vars["name"].(string)
	instance := gcp.FetchInstance(t, project, instanceName)

	// generate a ssh key pair
	keyPair := ssh.GenerateRSAKeyPair(t, 2048)

	// add the public ssh key to the compute engine metadata
	sshUsername := "terratest"
	publicKey := keyPair.PublicKey
	instance.AddSshKey(t, sshUsername, publicKey)

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "compute_engine_public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "terratest",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair,
		SshUserName: sshUsername,
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 10
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to public host %s", publicInstanceIP)

	// Run a simple echo command on the server
	expectedText := "Hello, World"
	command := fmt.Sprintf("echo -n '%s'", expectedText)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
	logger.Log(t, "--- PASS: testSSHToPublicHost")
}
