# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY SERVICE ACCOUNTS
# ---------------------------------------------------------------------------------------------------------------------

##### setup dbt bigquery service account to be attached to compute engine VM #####
resource "google_service_account" "service-account-dbt" {
  project      = var.project
  account_id   = var.account_id_dbt
  display_name = var.display_name_dbt
  description  = var.description_dbt
}

# hard-coded as this is specific to dbt enablement
# https://docs.getdbt.com/reference/warehouse-profiles/bigquery-profile/#required-permissions
locals {
  dbt_service_account_roles = concat(var.dbt_service_account_roles_to_add, [
    "roles/bigquery.user",
    "roles/bigquery.dataEditor"
  ])
}

resource "google_project_iam_binding" "dbt-iam-permissions" {
  project  = var.project
  for_each = toset(local.dbt_service_account_roles)
  role     = each.value

  members = [
    "serviceAccount:${google_service_account.service-account-dbt.email}",
  ]
}
