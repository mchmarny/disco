# Description: Creates a bucket for storing artifacts

# Creates a bucket for storing artifacts
resource "google_storage_bucket" "report_bucket" {
  name          = "${var.name}-${var.project_id}"
  location      = var.location
  storage_class = "STANDARD"
  force_destroy = true

  uniform_bucket_level_access = true
}

# Binds the runner service account to the bucket with bucket admin permissions
resource "google_storage_bucket_iam_binding" "report_bucket_binding" {
  bucket = google_storage_bucket.report_bucket.name
  role   = "roles/storage.admin"
  members = [
    "serviceAccount:${google_service_account.runner_service_account.email}",
  ]
}
