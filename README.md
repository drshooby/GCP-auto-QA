# GCP-auto-QA

automated QA test runner to learn several GCP services

# Authentication & Service Account Setup Guide for Automated QA Runner

This document explains how to configure authentication for this project, including:

- Creating a service account and its key file (`key.json`)
- Switching between service account credentials and user credentials
- Running Go tests/applications that require ID tokens
- Running Terraform locally using user credentials
- Best practices to avoid common pitfalls

---

## 1. Creating the Service Account and `key.json`

The service account enables your Go app to authenticate with Google Cloud, mint ID tokens, and call protected Cloud Run or Cloud Functions endpoints.

> **Security:** Treat `key.json` like a password. Never commit it to source control.

### Step-by-step:

1. **Create the service account**

```bash
gcloud iam service-accounts create automated-qa-runner \
  --description="Used for automated QA test runner" \
  --display-name="Automated QA Runner"
```

2. **Grant the service account necessary roles**

Replace YOUR_PROJECT_ID with your Google Cloud project ID.

```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:automated-qa-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.invoker"
```

3. **Generate and download the key file**

```bash
gcloud iam service-accounts keys create ~/key.json \
  --iam-account=automated-qa-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

Optionally, move it into your project directory:

```bash
mv ~/key.json ./key.json
```

4. **Add key.json to .gitignore**

```bash
# Ignore service account keys
key.json
```

## 2. Authentication Methods: When and How to Use

This project supports two authentication modes:

|       Mode       | Use for                        |                                         Setup                                          |              Notes               |
| :--------------: | :----------------------------- | :------------------------------------------------------------------------------------: | :------------------------------: |
| Service Account  | Go app/tests needing ID tokens |                     Set `GOOGLE_APPLICATION_CREDENTIALS=key.json`                      |  Required for minting ID tokens  |
| User Credentials | Terraform local usage          | Use `gcloud auth application-default login` and unset `GOOGLE_APPLICATION_CREDENTIALS` | Easy for Terraform, no ID tokens |

### Using Service Account Credentials (For Go app / tests)

```bash
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/key.json"  # or ./key.json if in project root
```

Run your Go tests or application:

```bash
go test -v ./...
# or
go run main.go
```

### Using User Credentials (For Terraform locally)

If you want to use your Google user credentials (via gcloud) instead of a service account:

1. Unset service account env var

```bash
unset GOOGLE_APPLICATION_CREDENTIALS
```

2. Login to gcloud application default credentials

```bash
gcloud auth application-default login
```

3. Run terraform

```bash
terraform init
terraform plan
terraform apply
```

Terraform will use the user credentials stored at:

```bash
~/.config/gcloud/application_default_credentials.json
```

## 3. Best Practices & Notes

- Never commit your key.json to source control.
- Use service accounts for automation, CI/CD, and anything requiring ID tokens.
- Use user credentials (gcloud auth application-default login) for local Terraform workflows.
- When running inside Cloud Run, Google’s default service account is used automatically (no need for key files).
- Always confirm that your service account has the correct IAM roles (e.g., roles/run.invoker for Cloud Run access).
