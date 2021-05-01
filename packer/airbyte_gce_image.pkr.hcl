variable "gcp_project_id" {
  type =  string
  default = "dbt-demos-sung"
}

variable "gcp_zone" {
  type =  string
  default = "us-central1-a"
}

variable "airbyte_build_script" {
  type =  string
  default = "airbyte_build.sh"
}

# debian hard coded as validated source image: https://docs.airbyte.io/deploying-airbyte/on-gcp-compute-engine

source "googlecompute" "airbyte-example" {
  project_id = var.gcp_project_id
  source_image = "debian-10-buster-v20210420"
  ssh_username = "packer"
  zone = var.gcp_zone
}

build {
    sources = ["sources.googlecompute.airbyte-example"]

    provisioner "shell" {
      script = var.airbyte_build_script
  }
}
