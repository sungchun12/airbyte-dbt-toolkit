
# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------
variable "credentials" {
  description = "path to service account json file"
  type        = string
  default     = "../service_account.json"
}

variable "project" {
  description = "name of your GCP project"
  type        = string
  default     = "dbt-demos-sung"
}

variable "location" {
  description = "default location of various GCP services"
  type        = string
  default     = "us-central1"
}

variable "subnetwork" {
  description = "default location of VPC subnet"
  type        = string
  default     = "projects/dbt-demos-sung/regions/us-central1/subnetworks/default"
}

variable "zone" {
  description = "default granular location typically for VMs"
  type        = string
  default     = "us-central1-a"
}

variable "service_account_email" {
  description = "Service account used for VMs"
  type        = string
  default     = "packer@dbt-demos-sung.iam.gserviceaccount.com"
}

variable "version_label" {
  description = "helpful label to version GCP resources per deployment"
  type        = string
  default     = "demo"
}

variable "image" {
  description = "OS image for compute engine instance"
  type        = string
  default     = "packer-1620060033"
}
