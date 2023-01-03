# Description: Cloud Run service

# List of roles that will be assigned to the runner service account
locals {
  # List of roles that will be assigned to the runner service account
  runner_roles = toset([
    "roles/artifactregistry.writer",
    "roles/bigquery.dataEditor",
    "roles/browser",
    "roles/containeranalysis.notes.attacher",
    "roles/iam.serviceAccountTokenCreator",
    "roles/monitoring.metricWriter",
    "roles/run.viewer",
    "roles/storage.objectCreator",
    "roles/storage.objectViewer",
    "roles/viewer",
  ])
}

# Service Account under which the Cloud Run services will run
resource "google_service_account" "runner_service_account" {
  account_id   = "${var.name}-run-sa"
  display_name = "Cloud Run service account for ${var.name}"
}

# Role binding
resource "google_project_iam_member" "runner_role_bindings" {
  for_each = local.runner_roles
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.runner_service_account.email}"
}

# Cloud Run service 
resource "google_cloud_run_service" "app" {
  name                       = var.name
  location                   = var.location
  project                    = var.project_id
  autogenerate_revision_name = true

  metadata {
    annotations = {
      "run.googleapis.com/ingress"     = "all"
      "run.googleapis.com/client-name" = "terraform"
    }
  }

  template {
    spec {
      containers {
        image = "${var.server_img}:${data.template_file.version.rendered}"

        ports {
          name           = "http1"
          container_port = 8080
        }
        resources {
          limits = {
            cpu    = "1000m"
            memory = "2Gi"
          }
        }
        env {
          name  = "ADDRESS"
          value = ":8080"
        }
        env {
          name  = "PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "GCS_BUCKET"
          value = google_storage_bucket.report_bucket.name
        }
      }

      container_concurrency = 80
      timeout_seconds       = 900
      service_account_name  = google_service_account.runner_service_account.email
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"         = "3"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# IAM member to grant access to the Cloud Run service
resource "google_cloud_run_service_iam_member" "app-access" {
  location = google_cloud_run_service.app.location
  project  = google_cloud_run_service.app.project
  service  = google_cloud_run_service.app.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.runner_service_account.email}"
}