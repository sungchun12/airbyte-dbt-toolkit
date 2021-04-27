source "googlecompute" "basic-example" {
  project_id = "dbt-demos-sung"
  source_image = "debian-9-stretch-v20200805"
  ssh_username = "packer"
  zone = "us-central1-a"
}

build {
  sources = ["sources.googlecompute.basic-example"]
}