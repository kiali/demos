variable "project_id" {
  type        = string
  description = "The GCP project that resources will be deployed in."
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "zone" {
  type    = string
  default = "us-central1-a"
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

variable "ssh_username" {
  type        = string
  description = "Username to be used when connecting via ssh. This is probably your google account username e.g. just <username> from <username>@company.com"
}

# Note: trial plans will probably hit the public ip address quota in GCP when adding more than 2 nodes.
variable "gke_num_nodes" {
  type        = number
  default     = 2
  description = "number of gke nodes"
}

variable "gke_machine_type" {
  type        = string
  default     = "n1-standard-2"
  description = "Machine type to use for gke nodes"
}

variable "deploy_bastion" {
  type        = bool
  default     = false
  description = "If set to true, a bastion host will be deployed providing access to the bookinfo vms."
}
