resource "google_cloudfunctions2_function_iam_member" "invoker_user" {
  location       = google_cloudfunctions2_function.function.location
  cloud_function = google_cloudfunctions2_function.function.name
  project        = var.project_id
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:automated-qa-runner@appspot.gserviceaccount.com"
}

resource "google_secret_manager_secret_iam_member" "slack_webhook_access" {
  secret_id = google_secret_manager_secret.slack_webhook.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_secret_manager_secret_version.slack_webhook_url_version]
}
