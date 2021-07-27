provider "google" {

  project = var.project_id
  region  = var.region
  zone    = var.zone
}

provider "google-beta" {

  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# This fetches a new token, which will expire in 1 hour.
data "google_client_config" "default" {}

provider "helm" {
  kubernetes {
    host = google_container_cluster.primary.endpoint

    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  }
}

provider "kubernetes" {
  host = google_container_cluster.primary.endpoint

  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

provider "kubectl" {
  host = google_container_cluster.primary.endpoint

  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  load_config_file       = false
}