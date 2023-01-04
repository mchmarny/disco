# Description: This file contains the resources required to federate a GitHub repository with a GCP project

# This is a list of roles that will be assigned to the GitHub federted user
locals {
  # List of roles that will be assigned to the GitHub federted user
  ci_roles = toset([
    "roles/artifactregistry.writer",
    "roles/bigquery.dataEditor",
    "roles/cloudkms.cryptoKeyEncrypterDecrypter",
    "roles/cloudkms.publicKeyViewer",
    "roles/cloudkms.signerVerifier",
    "roles/cloudkms.viewer",
    "roles/viewer",
  ])
}

# Service account to be used for federated auth to publish to GCR (existing)
resource "google_service_account" "github_actions_user" {
  count        = var.runtime_only ? 0 : 1
  account_id   = "${var.name}-github-actions-user"
  display_name = "Service Account impersonated in ${var.git_repo} GitHub Actions"
}

# IAM policy bindings to the service account resources created by GitHub identify
resource "google_project_iam_member" "ci_role_bindings" {
  for_each = local.ci_roles
  project  = var.project_id
  role     = var.runtime_only ? null : each.value
  member   = "serviceAccount:${google_service_account.runner_service_account.email}"
}

# Identiy pool for GitHub action based identity's access to Google Cloud resources
resource "google_iam_workload_identity_pool" "github_pool" {
  count                     = var.runtime_only ? 0 : 1
  provider                  = google-beta
  workload_identity_pool_id = "${var.name}-github-pool"
}

# Configuration for GitHub identiy provider
resource "google_iam_workload_identity_pool_provider" "github_provider" {
  count                              = var.runtime_only ? 0 : 1
  provider                           = google-beta
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool[count.index].workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.aud"        = "assertion.aud"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri        = "https://token.actions.githubusercontent.com"
    allowed_audiences = []
  }
}

# IAM policy bindings to the service account resources created by GitHub identify
resource "google_service_account_iam_member" "pool_impersonation" {
  count              = var.runtime_only ? 0 : 1
  provider           = google-beta
  service_account_id = google_service_account.github_actions_user[count.index].id
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool[count.index].name}/attribute.repository/${var.git_repo}"
}
