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
  source             = "./compute_engine"
  project            = var.project
  subnetwork_id      = module.networking.subnetwork
  bastion_host_email = module.service_accounts.service-account-bastion-host-email
  depends_on         = [module.api-enable-services]
}
