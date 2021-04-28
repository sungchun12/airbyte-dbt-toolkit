source "googlecompute" "airbyte-example" {
  project_id = "dbt-demos-sung"
  source_image = "debian-10-buster-v20210420"
  ssh_username = "packer"
  zone = "us-central1-a"
}

build {
    sources = ["sources.googlecompute.airbyte-example"]

    provisioner "shell" {
      script = "airbyte_build.sh"
  }
}
