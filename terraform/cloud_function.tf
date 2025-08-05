// https://cloud.google.com/functions/docs/tutorials/terraform

locals {
  project = var.project_id
}

resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name                        = "${random_id.default.hex}-gcf-source"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true // doesn't delete non-empty
}

data "archive_file" "default" {
  type        = "zip"
  output_path = "/tmp/function-source.zip"
  source_dir  = "${path.module}/../functions/status"
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.default.name
  source = data.archive_file.default.output_path
}

resource "google_cloudfunctions2_function" "function" {
  name        = "auto-qa-runner-function"
  description = "report status"
  location    = "us-west1"

  build_config {
    runtime     = "python311"
    entry_point = "status_get"
    source {
      storage_source {
        bucket = google_storage_bucket.default.name
        object = google_storage_bucket_object.object.name
      }
    }
  }


  service_config {
    max_instance_count = 1
    available_memory   = "128Mi"
    timeout_seconds    = 60
    secret_environment_variables {
      key        = "SLACK_WEBHOOK_URL"
      project_id = local.project
      secret     = google_secret_manager_secret.slack_webhook.secret_id
      version    = "latest"
    }
  }

  depends_on = [
    google_secret_manager_secret_version.slack_webhook_url_version,
    google_secret_manager_secret_iam_member.slack_webhook_access
  ]
}

resource "google_cloudfunctions2_function_iam_member" "invoker_user" {
  location       = google_cloudfunctions2_function.function.location
  cloud_function = google_cloudfunctions2_function.function.name
  project        = var.project_id
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:automated-qa-runner@appspot.gserviceaccount.com"
}
