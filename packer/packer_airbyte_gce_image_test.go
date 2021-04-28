package test

// install modules before running the test
// go get github.com/gruntwork-io/terratest/modules/docker
// go get github.com/stretchr/testify/assert@v1.4.0
// go get github.com/gruntwork-io/terratest/modules/packer@v0.34.0

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/stretchr/testify/assert"
)

func TestPackerAirbyteImage(t *testing.T) {
	packerOptions := &packer.Options{
		Template: "airbyte_gce_image.pkr.hcl",
	}

	packer.BuildArtifact(t, packerOptions)

	opts := &docker.RunOptions{Command: []string{"cat", "/test.txt"}}
	output := docker.Run(t, "gruntwork/packer-hello-world-example", opts)
	assert.Equal(t, "Hello, World!", output)
}
