output "service-account-dbt-email" {
  description = "The email for the dbt service account"
  value       = google_service_account.service-account-dbt.email
}
