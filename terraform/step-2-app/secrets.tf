resource "google_secret_manager_secret" "slack_webhook" {
  secret_id = "slack-webhook-url"

  labels = {
    label = "slack-results"
  }

  replication {
    auto {}
  }

  deletion_protection = false
}

resource "google_secret_manager_secret_version" "slack_webhook_url_version" {
  secret      = google_secret_manager_secret.slack_webhook.id
  secret_data = var.slack_webhook_url
}

data "google_project" "project" {
  project_id = var.project_id
}
