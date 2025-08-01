resource "google_artifact_registry_repository" "auto-qa-repo" {
  location      = "us-west1"
  repository_id = "auto-qa-repo"
  format        = "DOCKER"

  docker_config {
    immutable_tags = true
  }
}
