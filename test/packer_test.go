package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/packer"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
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

// This test builds the packer image and validates no breaking errors come through.
// The test is broken into "stages" so you can skip stages by setting environment variables (e.g.,
// skip stage "build_image" by setting the environment variable "SKIP_build_image=true"), which speeds up iteration when
// running this test over and over again locally.
func TestPackerAirbyteDemo(t *testing.T) {
	t.Parallel()

	// The folder where we have our packer code
	workingDir := "../packer/"

	// Get the Project Id to use
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// At the end of the test, delete the Image
	defer test_structure.RunTestStage(t, "cleanup_image", func() {
		imageName := test_structure.LoadString(t, workingDir, "imageName")
		deleteImage(t, projectID, imageName)
	})

	// Build the Image for the web app
	test_structure.RunTestStage(t, "build_image", func() {
		buildImage(t, projectID, workingDir)
	})
}

// Build the Packer Image
func buildImage(t *testing.T, projectID string, workingDir string) {

	// Pick a random GCP zone to test in. This helps ensure your code works in all regions.
	// On October 22, 2018, GCP launched the asia-east2 region, which promptly failed all our tests, so blacklist asia-east2.
	zone := gcp.GetRandomZone(t, projectID, nil, nil, []string{"asia-east2"})
	airbyte_build_script := workingDir + "airbyte_build.sh"

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../packer/airbyte_gce_image.pkr.hcl",

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
	test_structure.SavePackerOptions(t, workingDir, packerOptions)

	// Make sure the Packer build completes successfully
	imageName := packer.BuildArtifact(t, packerOptions)

	// Save the imageName as a string so future test stages can use them
	test_structure.SaveString(t, workingDir, "imageName", imageName)
}

// Delete the Packer Image
func deleteImage(t *testing.T, projectID string, imageName string) {
	// Load the Image ID saved by the earlier build_image stage
	image := gcp.FetchImage(t, projectID, imageName)
	image.DeleteImage(t)
}
