package test

// go get github.com/gruntwork-io/terratest/modules/terraform@v0.34.0

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAirbyteDemo(t *testing.T) {
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform/",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "compute_engine_id")
	assert.Equal(t, "projects/dbt-demos-sung/zones/us-central1-a/instances/airbyte-demo", output)
}
