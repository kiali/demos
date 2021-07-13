
# Firewall
resource "google_compute_firewall" "allow-ssh-bastion" {
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }

  description   = "Allow SSH from anywhere for bastion"
  direction     = "INGRESS"
  disabled      = "false"
  name          = "allow-ssh-bastion"
  network       = google_compute_network.vpc.name
  priority      = "1100"
  project       = var.project_id
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["bastion"]
}

# Firewall
resource "google_compute_firewall" "allow-all-tcp-udp" {
  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
    ports    = ["53"]
  }

  description   = "Allow TCP from anywhere for bookinfo vms"
  direction     = "INGRESS"
  disabled      = "false"
  name          = "allow-all-tcp-udp-bookinfo"
  network       = google_compute_network.vpc.name
  priority      = "1100"
  project       = var.project_id
  source_ranges = ["0.0.0.0/0"]
  target_tags   = local.common_vm_tags
}

# Cloud DNS
# Attempts to override the /etc/hosts file in GCP were unsuccessful so instead
# we setup a private zone to resolve the istiod.<ns>.svc hostname on the VMs.
resource "google_dns_managed_zone" "istio-system-zone" {
  name        = "istio-system"
  dns_name    = "istio-system.svc."
  description = "Resolves the istio namespace for the bookinfo VMs"

  visibility = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.vpc.id
    }
  }

  depends_on = [
    module.project-services
  ]
}

resource "google_dns_record_set" "istiod_a_record" {
  managed_zone = google_dns_managed_zone.istio-system-zone.name
  name         = "istiod.istio-system.svc."
  type         = "A"
  rrdatas      = ["${data.kubernetes_service.eastwest-gateway-service.status.0.load_balancer.0.ingress.0.ip}"]
  ttl          = 300

  depends_on = [
    module.project-services
  ]
}

resource "kubectl_manifest" "vm_namespace" {
  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: bookinfo
  labels:
    istio-injection: enabled
YAML

}

resource "kubectl_manifest" "vm_service_account" {
  yaml_body = <<YAML
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm-service-account
  namespace: "${kubectl_manifest.vm_namespace.name}"
YAML

}

resource "local_file" "workload_entry_template" {
  for_each = local.workloads

  filename = "${path.module}/manifests/${each.key}_workloadgroup.yaml"
  content  = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: WorkloadGroup
metadata:
  name: ${each.key}
  namespace: ${kubectl_manifest.vm_namespace.name}
spec:
  metadata:
    labels:
      app: ${each.value.app}
      version: ${each.value.version}
      workload-type: vm
  template:
    serviceAccount: ${kubectl_manifest.vm_service_account.name}
    network: "${google_compute_network.vpc.name}"
    labels:
      app: ${each.value.app}
      version: ${each.value.version}
      workload-type: vm

YAML

}

resource "time_sleep" "wait_90_seconds_for_istio_config" {
  create_duration = "90s"

  depends_on = [
    kubectl_manifest.istio_config
  ]
}

resource "null_resource" "login_to_gke_cluster" {
  provisioner "local-exec" {
    # This will setup the local kube config to the talk to the gke cluster.
    command     = "gcloud container clusters get-credentials ${google_container_cluster.primary.name} --region ${var.region} --project ${var.project_id}"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [
    google_container_cluster.primary,
    google_container_node_pool.primary_nodes
  ]
}

resource "null_resource" "generate_workload_entry_files" {
  for_each = local_file.workload_entry_template

  triggers = {
    workload_entry_template_content = "${each.key}: ${each.value.content}"
    // TODO: Is there a trigger that doesn't cause a cycle we can use to run this whenever the vm changes?
  }

  provisioner "local-exec" {
    command     = "istioctl x workload entry configure -f ${path.module}/manifests/${each.key}_workloadgroup.yaml -o ${path.module}/data/${each.key} --clusterID ${google_container_cluster.primary.name}"
    interpreter = ["/bin/bash", "-c"]
  }

  depends_on = [
    # istioctl needs to talk to the cluster and will timeout or fail if the kubeconfig is not setup properly.
    null_resource.login_to_gke_cluster,
    # Giving the istio operator a chance to create istiod instance. Probably a better way to do this than sleeping.
    time_sleep.wait_90_seconds_for_istio_config,
  ]
}

resource "kubectl_manifest" "kiali_traffic_generator" {
  yaml_body = <<YAML
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: kiali-traffic-generator
  labels:
    app: kiali-traffic-generator
    kiali-test: traffic-generator
    istio-injection: disabled
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kiali-traffic-generator
      kiali-test: traffic-generator
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: kiali-traffic-generator
        kiali-test: traffic-generator
    spec:
      containers:
      - image: kiali/kiali-test-mesh-traffic-generator:latest
        name: kiali-traffic-generator
        imagePullPolicy: IfNotPresent
        env:
        - name: DURATION
          value: "0s"
        - name: ROUTE
          value: "http://${data.kubernetes_service.ingress-gateway-service.status.0.load_balancer.0.ingress.0.ip}/productpage"
        - name: RATE
          value: "1"

YAML

  depends_on = [
    helm_release.istio_base,
    # Giving the istio operator a chance to create istiod instance. Probably a better way to do this than sleeping.
    time_sleep.wait_90_seconds_for_istio_config,
  ]
}

resource "kubectl_manifest" "bookinfo_deployment" {
  for_each = var.deploy_in_kube ? local.workloads : {}

  wait_for_rollout = true

  yaml_body = <<YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "${each.key}"
  namespace: "${kubectl_manifest.vm_namespace.name}"
  labels:
    app: "${each.value.app}"
    version: "${each.value.version}"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: "${each.value.app}"
      version: "${each.value.version}"
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: "${each.value.app}"
        version: "${each.value.version}"
    spec:
      serviceAccount: "${kubectl_manifest.vm_service_account.name}"
      containers:
      - image: "istio/examples-bookinfo-${each.value.app}-${each.value.version}:latest"
        name: "${each.value.app}"
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080

YAML

  # Really the istio controlplane just needs to be online before we
  # create the deployments so that the sidecar injection will happen
  # but waiting until after the vm workloads gives more time for
  # the controlplane to be ready.
  depends_on = [
    helm_release.istio_operator,
    google_compute_instance.workload_vm
  ]
}

# Any address can be used in the addresses field as long as it's not
# localhost as the envoy will end up sending the request to the actual
# endpoint and not the address that was resolved by DNS. We need a value
# there so that the app can resolve something. Otherwise productpage
# would not resolve to anything from the proxy and the request would
# get forwarded on to an upstream DNS server where it definitely
# wouldn't resolve.
resource "kubectl_manifest" "service_entry" {
  for_each = toset([for workload in local.workloads : workload.app])

  yaml_body = <<YAML
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: "${each.key}"
  namespace: "${kubectl_manifest.vm_namespace.name}"
spec:
  addresses:
    - 8.8.8.8
  hosts:
    - "${each.key}"
    - "${each.key}.${kubectl_manifest.vm_namespace.name}.svc.cluster.local"
  location: MESH_INTERNAL
  resolution: STATIC
  ports:
  - number: 9080
    name: http
    protocol: HTTP
    targetPort: 9080
  workloadSelector:
    labels:
      app: "${each.key}"

YAML

  depends_on = [
    helm_release.istio_base
  ]
}


resource "kubectl_manifest" "workload_entry" {
  for_each = local.workloads_with_vm_info
  # Without this, updating the vm does not update the service entry during planning phase.
  force_new = true

  yaml_body = <<YAML
apiVersion: networking.istio.io/v1beta1
kind: WorkloadEntry
metadata:
  name: "${each.value.name}"
  namespace: "${kubectl_manifest.vm_namespace.name}"
spec:
  serviceAccount: "${kubectl_manifest.vm_service_account.name}"
  address: "${each.value.network_interface.0.network_ip}"
  network: "${google_compute_network.vpc.name}"
  labels:
    app: "${each.value.app}"
    version: "${each.value.version}"
    workload-type: vm

YAML

  depends_on = [
    helm_release.istio_base,
    null_resource.generate_workload_entry_files,
  ]
}

resource "kubectl_manifest" "product_page_vs" {
  yaml_body = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: productpage
  namespace: "${kubectl_manifest.vm_namespace.name}"
spec:
  hosts:
    - "*"
  gateways:
    - ${kubectl_manifest.vm_namespace.name}/bookinfo-gateway
  http:
    - match:
        - uri:
            exact: /productpage
        - uri:
            exact: /login
        - uri:
            exact: /logout
        - uri:
            prefix: /api/v1/products
        - uri:
            prefix: /static
      route:
        - destination:
            host: productpage.${kubectl_manifest.vm_namespace.name}.svc.cluster.local
YAML

  depends_on = [
    helm_release.istio_base
  ]
}

resource "kubectl_manifest" "workload_vs" {
  for_each = setsubtract(local.apps, ["productpage"]) # Productpage gets its own special VS

  yaml_body = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: "${each.key}"
  namespace: "${kubectl_manifest.vm_namespace.name}"
spec:
  hosts:
    - "${each.key}"
  http:
    - route:
      - destination:
          host: "${each.key}.${kubectl_manifest.vm_namespace.name}.svc.cluster.local"
YAML

  depends_on = [
    helm_release.istio_base
  ]
}

resource "kubectl_manifest" "ingress_gateway" {
  yaml_body = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: bookinfo-gateway
  namespace: "${kubectl_manifest.vm_namespace.name}"
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
YAML 

  depends_on = [
    helm_release.istio_base
  ]
}

data "local_file" "root_cert" {
  filename = "${path.module}/data/${each.key}/root-cert.pem"

  for_each = local.workloads

  depends_on = [
    null_resource.generate_workload_entry_files,
  ]
}

data "local_file" "cluster_env" {
  filename = "data/${each.key}/cluster.env"

  for_each = local.workloads

  depends_on = [
    null_resource.generate_workload_entry_files,
  ]
}

data "local_file" "istio_token" {
  filename = "data/${each.key}/istio-token"

  for_each = local.workloads

  depends_on = [
    null_resource.generate_workload_entry_files,
  ]
}

data "local_file" "mesh" {
  filename = "data/${each.key}/mesh.yaml"

  for_each = local.workloads

  depends_on = [
    null_resource.generate_workload_entry_files,
  ]
}

# VM
resource "google_compute_instance" "workload_vm" {
  for_each = local.workload_startup_scripts

  allow_stopping_for_update = true
  boot_disk {
    auto_delete = "true"
    device_name = "boot-disk"
    initialize_params {
      image = "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-10-buster-v20210609"
      size  = "10"
      type  = "pd-standard"
    }

    mode = "READ_WRITE"
  }

  can_ip_forward      = "false"
  deletion_protection = "false"
  enable_display      = "false"
  machine_type        = "n1-standard-1"

  metadata = {
    ssh-keys = "${var.ssh_username}:${local.ssh_key}"
  }

  metadata_startup_script = each.value

  name = each.key

  network_interface {
    network    = google_compute_network.vpc.name
    subnetwork = google_compute_subnetwork.subnet.name
  }

  project = var.project_id

  scheduling {
    automatic_restart   = "true"
    on_host_maintenance = "MIGRATE"
    preemptible         = "false"
  }

  service_account {
    scopes = ["cloud-platform"]
  }

  tags = concat(local.common_vm_tags, ["${each.key}"])

  zone = var.zone

  depends_on = [
    module.project-services,
    null_resource.generate_workload_entry_files,
  ]
}

# Cloud NAT
module "cloud_router" {
  source  = "terraform-google-modules/cloud-router/google"
  version = "~> 0.4"
  project = var.project_id
  name    = "bookinfo-cloud-router"
  network = google_compute_network.vpc.name
  region  = var.region

  nats = [{
    name = "bookinfo-nat-gateway"
  }]
}


# Bastion
resource "google_compute_instance" "bastion" {
  count                     = var.deploy_bastion ? 1 : 0
  allow_stopping_for_update = true
  boot_disk {
    auto_delete = "true"
    device_name = "boot-disk"
    initialize_params {
      image = "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-10-buster-v20210609"
      size  = "10"
      type  = "pd-standard"
    }

    mode = "READ_WRITE"
  }

  can_ip_forward      = "false"
  deletion_protection = "false"
  enable_display      = "false"
  machine_type        = "f1-micro"

  metadata = {
    ssh-keys = "${var.ssh_username}:${local.ssh_key}"
  }

  name = "bastion"

  network_interface {
    access_config {}

    network    = google_compute_network.vpc.name
    subnetwork = google_compute_subnetwork.subnet.name
  }

  project = var.project_id

  scheduling {
    automatic_restart   = "true"
    on_host_maintenance = "MIGRATE"
    preemptible         = "false"
  }

  service_account {
    scopes = ["cloud-platform"]
  }

  tags = ["bastion"]

  zone = var.zone

  depends_on = [
    module.project-services,
    null_resource.generate_workload_entry_files,
  ]
}