variable "project" {
  type =  string
  default = "dbt-demos-sung"
}

variable "zone" {
  type =  string
  default = "us-central1-a"
}

variable "airbyte_build_script" {
  type =  string
  default = "airbyte_build.sh"
}

# debian hard coded as validated source image: https://docs.airbyte.io/deploying-airbyte/on-gcp-compute-engine
variable "airbyte_source_image" {
  type =  string
  default = "debian-10-buster-v20210420"
}

variable "airbyte_ssh_username" {
  type =  string
  default = "packer"
}

source "googlecompute" "airbyte-example" {
  project_id = var.project
  source_image = var.airbyte_source_image
  ssh_username = var.airbyte_ssh_username
  zone = var.zone
}

build {
    sources = ["sources.googlecompute.airbyte-example"]

    provisioner "shell" {
      script = var.airbyte_build_script
  }
}
