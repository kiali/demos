# Terraform GCP VM Workload Example

This folder contains an example of a bookinfo istio deployment running on VMs in GCP. Its primary purpose is to demonstrate and test kiali features with VM workloads and not to be a reference for how to integrate VM workloads with Istio. It will deploy the following:

- GKE cluster
- Istio within the cluster
- VMs to run bookinfo
- VPC - single network for the vms and kube nodes

Requests to the bookinfo services should be spread between the vms and the kubernetes deployments for each app.

## Requirements

- [Terraform 1.x](https://learn.hashicorp.com/tutorials/terraform/install-cli#install-terraform)
- [istioctl](https://istio.io/latest/docs/setup/getting-started/#download) in your PATH
- [gcloud-cli](https://cloud.google.com/sdk/docs/install)

## Setup

1. Create a GCP project 
2. Login to gcloud with:

```bash
gcloud auth application-default login
```

**Note**: Your account should also have billing enabled for your project. This project can be deployed using a free trial account without needing to upgrade but either way you'll need to ensure you've activated billing first in the GCP console.

## Run

```bash
terraform init
terraform apply
```

You can either enter values for your variables through the command line prompt or create a `terraform.tfvars` file to store them in:

```
project_id = ""
region     = ""
ssh_key = ""
zone = ""
ssh_username = ""
```

## Access

After the terraform runs successfully it will output access information for various resources.

### Kubernetes - Kubectl

Your kube config's current context will already be changed to the gke cluster after a successful deployment but if you need to re-authenticate run:

```bash
gcloud container clusters get-credentials <kubernetes_cluster_name> --region <region>
```

### VM

None of the workload VMs have external ips but they can be accessed via an optional bastion.

Set the `deploy_bastion = true` terraform variable and re-run `terraform apply`. You can then jump through to any of the workload vms through the bastion by running

```bash
ssh -J <ssh_username>@<bastion_ip> <ssh_username>@<internal_vm_ip>
```

### Kiali

```bash
kubectl port-forward service/kiali -n istio-system 20001:20001
```

then naviage to https://localhost:20001 in your browser.
