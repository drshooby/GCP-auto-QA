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

# Get the project data to construct the default compute service account
data "google_project" "project" {
  project_id = var.project_id
}

# IAM binding to allow the Cloud Function to access the secret
# Cloud Functions v2 uses the default compute service account by default
resource "google_secret_manager_secret_iam_member" "slack_webhook_access" {
  secret_id = google_secret_manager_secret.slack_webhook.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_secret_manager_secret_version.slack_webhook_url_version]
}
