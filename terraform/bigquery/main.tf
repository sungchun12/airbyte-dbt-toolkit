resource "google_bigquery_dataset" "airbyte-dataset" {
  project                     = var.project
  dataset_id                  = var.airbyte_dataset_id
  friendly_name               = var.airbyte_dataset_friendly_name
  description                 = var.airbyte_dataset_description
  location                    = var.airbyte_dataset_location
  default_table_expiration_ms = var.airbyte_dataset_table_expiration

  labels = var.airbyte_dataset_labels

}
