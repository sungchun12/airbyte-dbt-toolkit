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

  #   access {
  #     role          = "OWNER"
  #     user_by_email = google_service_account.bqowner.email
  #   }

}

# resource "google_service_account" "bqowner" {
#   account_id = "bqowner"
# }
