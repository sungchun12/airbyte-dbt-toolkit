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
variable "account_id_dbt" {
  description = "Service account, account id for dbt"
  type        = string
  default     = "dbt-service-account"
}

variable "display_name_dbt" {
  description = "Display name for the dbt service account"
  type        = string
  default     = "dbt Service Account for accessing BigQuery"
}

variable "description_dbt" {
  description = "Description for the dbtservice account"
  type        = string
  default     = "Service account for dev purposes to access BigQuery with dbt commands"
}
variable "dbt_service_account_roles_to_add" {
  description = "Additional roles to be added to the service account."
  type        = list(string)
  default     = []
}
