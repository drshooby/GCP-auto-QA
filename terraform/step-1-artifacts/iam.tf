# resource "google_project_iam_member" "artifact_registry_writer" {
#   project = var.project_id
#   role    = "roles/artifactregistry.writer"
#   member  = "serviceAccount:${google_service_account.artifact_service_account.email}"
# }

# resource "google_project_iam_member" "artifact_registry_reader" {
#   project = var.project_id
#   role    = "roles/artifactregistry.reader"
#   member  = "serviceAccount:${google_service_account.artifact_service_account.email}"
# }

# resource "google_service_account_iam_member" "allow_impersonation" {
#   service_account_id = google_service_account.artifact_service_account.name
#   role               = "roles/iam.serviceAccountTokenCreator"
#   member             = "user:${var.user-email}"
# }
