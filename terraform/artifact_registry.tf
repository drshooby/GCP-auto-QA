resource "google_artifact_registry_repository" "auto-qa-repo" {
  location      = "us-west1"
  repository_id = "auto-qa-repo"
  format        = "DOCKER"

  docker_config {
    immutable_tags = false # doesn't let me delete images without this
  }
}

resource "google_service_account" "artifact_service_account" {
  account_id   = "runner-artifacts-service"
  display_name = "Service Account for automated QA runner artifact pushes"
  project      = "automated-qa-runner"
}

resource "google_project_iam_member" "artifact_registry_writer" {
  project = "automated-qa-runner"
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.artifact_service_account.email}"
}

resource "google_project_iam_member" "artifact_registry_reader" {
  project = "automated-qa-runner"
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.artifact_service_account.email}"
}

resource "google_service_account_iam_member" "allow_impersonation" {
  service_account_id = google_service_account.artifact_service_account.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "user:${var.user-email}"
}
