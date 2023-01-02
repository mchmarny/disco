# Description: Cloud Scheduler job to update user data

# Cloud Scheduler job to update user data
resource "google_cloud_scheduler_job" "disco_refresh_job" {
  name             = "${var.name}-disco-refresh-job"
  description      = "Invokes Cloud Run service to run discovery"
  schedule         = "42 */5 * * 0"
  time_zone        = "America/Los_Angeles"
  attempt_deadline = "900s"
  region           = var.location

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "GET"
    uri         = "${google_cloud_run_service.app.status[0].url}/disco"

    oidc_token {
      service_account_email = google_service_account.runner_service_account.email
      audience              = "${google_cloud_run_service.app.status[0].url}/disco"
    }
  }
}
