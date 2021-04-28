# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY COMPUTE ENGINE BASTION HOST TO RUN COMMANDS AGAINST THE CLOUD COMPOSER KUBERNETES CLUSTER
# ---------------------------------------------------------------------------------------------------------------------

# Minimal Setup
resource "google_compute_instance" "airbyte-demo" {
  project      = var.project
  name         = var.name
  machine_type = var.machine_type
  zone         = var.zone

  tags = var.tags

  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = var.image
    }
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = true
    enable_vtpm                 = true
  }

  scratch_disk {
    interface = var.interface
  }

  network_interface {
    subnetwork = var.subnetwork_id

    access_config {
      // Ephemeral IP
    }
  }

  metadata = var.metadata

  metadata_startup_script = <<SCRIPT
  ${file("${path.module}/${var.metadata_startup_script}")}
  SCRIPT

  service_account {
    email  = var.service_account_email
    scopes = var.scopes
  }
}
