output "service-account-dbt-email" {
  description = "The email for the dbt service account"
  value       = google_service_account.service-account-dbt.email
}

output "dbt-iam-permissions-list" {
  description = "The list of dbt iam permissions"
  value       = [for x in google_project_iam_binding.dbt-iam-permissions : x.role]
}
