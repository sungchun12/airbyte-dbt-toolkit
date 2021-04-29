# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# DEPLOY A COMPREHENSIVE CLOUD COMPOSER ENVIRONMENT WITH SUPPORTING INFRASTRUCTURE AND NETWORKING
# This module serves as the main interface for all the supporting modules
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# ---------------------------------------------------------------------------------------------------------------------
# IMPORT MODULES
# This root module imports and passes through project wide variables
# Detailed default variables contained within respective module directory
# -------------------------------------------------------------------------------------------------------------------


module "compute_engine" {
  source                = "./compute_engine"
  project               = var.project
  subnetwork_id         = var.subnetwork
  service_account_email = var.service_account_email
}

module "service_accounts" {
  source  = "./service_accounts"
  project = var.project
}

module "bigquery" {
  source  = "./bigquery"
  project = var.project
}
