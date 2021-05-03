# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------
variable "project" {
  description = "project where terraform will setup these services"
  type        = string
}

variable "subnetwork_id" {
  description = "Subnetwork where bastion host will be setup, defaults to the same one as cloud composer"
  type        = string
}

variable "service_account_email" {
  description = "Service account email with cloud composer permissions"
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL MODULE PARAMETERS
# These variables have defaults, but may be overridden by the operator
# ---------------------------------------------------------------------------------------------------------------------

variable "name" {
  description = "compute engine instance name"
  type        = string
  default     = "airbyte-demo"
}

variable "machine_type" {
  description = "compute engine machine type"
  type        = string
  default     = "n1-standard-2"
}

variable "zone" {
  description = "zone where the compute engine instace is setup"
  type        = string
  default     = "us-central1-a"
}

variable "tags" {
  description = "tags to associate with the compute engine instance"
  type        = list(string)
  default     = ["dev"]
}

variable "image" {
  description = "OS image for compute engine instance"
  type        = string
  default     = null
}


variable "interface" {
  description = "Local SSD disk"
  type        = string
  default     = "SCSI"
}


variable "metadata" {
  description = "Extra metadata configs to set"
  type        = map(string)
  default = {
    block-project-ssh-keys = false
  }
}


variable "metadata_startup_script" {
  description = "Metadata startup script after the instance is setup"
  type        = string
  default     = "airbyte_deploy.sh"
}

variable "scopes" {
  description = "Determine which services this compute engine instance is allowed to interact with"
  type        = list(string)
  default     = ["https://www.googleapis.com/auth/cloud-platform"]
}

