output "airbyte_dataset_id" {
  description = "The exact id of the bigquery dataset id"
  value       = google_bigquery_dataset.airbyte-dataset.id
}
