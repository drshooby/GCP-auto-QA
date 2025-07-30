// https://cloud.google.com/functions/docs/tutorials/terraform

resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name                        = "${random_id.default.hex}-gcf-source"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = true // TODO! why doesn't this delete the bucket
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

resource "google_cloudfunctions_function" "function" {
  name        = "auto-qa-runner-function"
  region      = "us-west1"
  description = "report status"
  runtime     = "python311"

  max_instances = 1

  source_archive_bucket = google_storage_bucket.default.name
  source_archive_object = google_storage_bucket_object.object.name

  available_memory_mb          = 128
  trigger_http                 = true
  https_trigger_security_level = "SECURE_ALWAYS"
  timeout                      = 60
  entry_point                  = "status_get"
}
