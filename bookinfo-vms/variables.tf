variable "project_id" {
  type        = string
  description = "The GCP project that resources will be deployed in."
}

variable "ssh_username" {
  type        = string
  description = "Username to be used when connecting via ssh. This is probably your google account username e.g. just <username> from <username>@company.com"
}

variable "region" {
  type        = string
  description = "Region where GCP resources will be deployed."
  default     = "us-central1"
}

variable "zone" {
  type        = string
  description = "Zone where GCP resources will be deployed."
  default     = "us-central1-a"
}

variable "ssh_key" {
  type        = string
  description = "Your public ssh key that will be used to access the VMs remotely. Takes precedence over ssh_key_filepath."
  default     = ""
}

variable "ssh_key_filepath" {
  type        = string
  description = "The path to your public ssh key that will be used to access the VMs remotely."
  default     = ""
}

# Note: trial plans will probably hit the public ip address quota in GCP when adding more than 2 nodes.
variable "gke_num_nodes" {
  type        = number
  description = "number of gke nodes"
  default     = 2
}

variable "gke_machine_type" {
  type        = string
  description = "Machine type to use for gke nodes"
  default     = "n1-standard-2"
}

variable "deploy_bastion" {
  type        = bool
  description = "If set to true, a bastion host will be deployed providing access to the bookinfo vms."
  default     = false
}

variable "deploy_in_kube" {
  type        = bool
  description = "If true, deploys the bookinfo apps inside the kubernetes cluster in addition to the vms."
  default     = true
}