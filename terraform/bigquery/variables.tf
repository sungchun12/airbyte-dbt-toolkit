# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------
variable "project" {
  description = "project where terraform will setup these services"
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL MODULE PARAMETERS
# These variables have defaults, but may be overridden by the operator
# ---------------------------------------------------------------------------------------------------------------------
variable "airbyte_dataset_id" {
  description = "Dataset ID for where airbyte will land the ingested data"
  type        = string
  default     = "airbyte_dataset"
}

variable "airbyte_dataset_friendly_name" {
  description = "Friendly name dataset ID for where airbyte will land the ingested data"
  type        = string
  default     = "airbyte_dataset"
}

variable "airbyte_dataset_description" {
  description = "airbyte_dataset_description"
  type        = string
  default     = "This is the main dataset for airbyte to point for EL destinations"
}

variable "airbyte_dataset_location" {
  description = "Location of the BigQuery dataset"
  type        = string
  default     = "US"
}

variable "airbyte_dataset_table_expiration" {
  description = "Default table expiration within the BigQuery dataset"
  type        = string
  default     = null
}

variable "airbyte_dataset_labels" {
  description = "Labels for the BigQuery dataset"
  type        = map(any)
  default     = { env = "dev" }
}


