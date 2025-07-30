#!/bin/bash
set -euo pipefail

STATUS=$1

TOKEN=$(gcloud auth print-identity-token)
curl -H "Authorization: Bearer $TOKEN" "$(cd ../terraform && terraform output -raw cloud_function_invoke_url)?status=${STATUS}"