resource "google_cloud_run_v2_job" "test_runner" {
  name                = "qa-test-runner"
  location            = "us-west1"
  deletion_protection = false

  template {
    template {
      containers {
        image = "us-west1-docker.pkg.dev/automated-qa-runner/auto-qa-repo/gin-math:latest"
        resources {
          limits = {
            cpu    = "2"
            memory = "1024Mi"
          }
        }
      }
    }
  }

  depends_on = [google_cloudfunctions2_function.function]
}
