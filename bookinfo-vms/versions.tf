terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.72.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.2.0"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}
