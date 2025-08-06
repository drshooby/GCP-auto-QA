resource "google_artifact_registry_repository" "auto-qa-repo" {
  location      = "us-west1"
  repository_id = "auto-qa-repo"
  format        = "DOCKER"
  project       = "automated-qa-runner"

  docker_config {
    immutable_tags = false # doesn't let me delete images without this
  }
}

resource "google_service_account" "artifact_service_account" {
  account_id   = "runner-artifacts-service"
  display_name = "Service Account for automated QA runner artifact pushes"
  project      = "automated-qa-runner"
}
