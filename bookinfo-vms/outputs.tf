output "region" {
  value       = var.region
  description = "GCloud Region"
}

output "project_id" {
  value       = var.project_id
  description = "GCloud Project ID"
}

output "kubernetes_cluster_name" {
  value       = google_container_cluster.primary.name
  description = "GKE Cluster Name"
}

output "kubernetes_cluster_host" {
  value       = google_container_cluster.primary.endpoint
  description = "GKE Cluster Host"
}

output "bastion" {
  value       = length(google_compute_instance.bastion) == 1 ? "${var.ssh_username}@${google_compute_instance.bastion[0].network_interface.0.access_config.0.nat_ip}" : "No bastion host. Set 'deploy_bastion' variable to deploy bastion."
  description = "Bastion host"
}

output "private_ips" {
  value       = join(" ", [for vm in google_compute_instance.workload_vm : "${vm.name}: ${vm.network_interface.0.network_ip}"])
  description = "private vm ips"
}

output "bookinfo_url" {
  value       = "http://${data.kubernetes_service.ingress-gateway-service.status.0.load_balancer.0.ingress.0.ip}/productpage"
  description = "URL to access the bookinfo productpage"
}
