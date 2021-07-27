# Enable GCP APIs needed by this project
module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "11.0.0"

  project_id                  = var.project_id
  disable_services_on_destroy = false
  disable_dependent_services  = false

  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "dns.googleapis.com",
    "iam.googleapis.com",
  ]
}

# VPC
resource "google_compute_network" "vpc" {
  name                    = "${var.project_id}-vpc"
  auto_create_subnetworks = "false"

  depends_on = [
    module.project-services
  ]
}

# Subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "${var.project_id}-subnet"
  region        = var.region
  network       = google_compute_network.vpc.name
  ip_cidr_range = "10.10.0.0/24"

  depends_on = [
    module.project-services
  ]
}

# GKE cluster
resource "google_container_cluster" "primary" {
  name     = "${var.project_id}-gke"
  location = var.region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  depends_on = [
    module.project-services
  ]
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "${google_container_cluster.primary.name}-node-pool"
  location   = var.region
  cluster    = google_container_cluster.primary.name
  node_count = var.gke_num_nodes

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    labels = {
      env = var.project_id
    }

    preemptible  = true
    machine_type = var.gke_machine_type
    tags         = ["gke-node", "${var.project_id}-gke"]
    metadata = {
      disable-legacy-endpoints = "true"
    }

  }

  depends_on = [
    module.project-services
  ]
}