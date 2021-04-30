resource "google_bigquery_dataset" "airbyte-dataset" {
  project                     = var.project
  dataset_id                  = "airbyte_dataset"
  friendly_name               = "airbyte_dataset"
  description                 = "This is the main dataset for airbyte to point for EL destinations"
  location                    = "US"
  default_table_expiration_ms = null

  labels = {
    env = "dev"
  }

}
