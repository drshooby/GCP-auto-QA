# resource "google_cloud_run_v2_service" "default" {
#   name                = "test-box"
#   location            = "us-west1"
#   deletion_protection = false

#   template {
#     scaling {
#       max_instance_count = 1
#     }
#     containers {
#       image = ""

#       env {
#         name = "GCP_AUTH_TOKEN"
#         value_source {
#           secret_key_ref {
#             secret  = google_secret_manager_secret.secret_one.secret_id
#             version = "lastest"
#           }
#         }
#       }

#       env {
#         name = "CLOUD_FN_URL"
#         value_source {
#           secret_key_ref {
#             secret  = google_secret_manager_secret.secret_two.secret_id
#             version = "latest"
#           }
#         }
#       }
#     }
#   }
# }
