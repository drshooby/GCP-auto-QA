terraform {
  required_version = "~> 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

// https://registry.terraform.io/providers/hashicorp/google/latest/docs
provider "google" {
  project = "automated-qa-runner"
  region  = "us-west1"
  // zone = ""
}
