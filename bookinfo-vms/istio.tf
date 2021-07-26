// When a new version is released and old charts dir exists
// it doesn't get cleaned up by tf destroy and this will throw an
// error because it'll pick up old version of charts during plan
// but new version will get downloaded after. Deleting the charts
// at destroy ensures that the version info won't be inconsistent
// during the plan process.
resource "null_resource" "delete_istio_charts" {
  provisioner "local-exec" {
    command     = "rm -rf ${path.module}/charts"
    interpreter = ["/bin/bash", "-c"]
    when        = destroy
  }
}

resource "null_resource" "download_istio_charts" {
  provisioner "local-exec" {
    command     = "scripts/download-istio.sh ${path.module}"
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "helm_release" "istio_base" {
  atomic           = true
  chart            = "./charts/base"
  create_namespace = true
  name             = "istio-base"
  namespace        = "istio-system"
  wait_for_jobs    = true

  depends_on = [
    google_container_cluster.primary,
    google_container_node_pool.primary_nodes,
    null_resource.download_istio_charts,
  ]
}

// The helm release doesn't like the fact that the namespace is included
// as a template in the istio operator chart so creating the namespace
// first outside of the istio operator chart and the namespace template
// has been removed from that chart.
resource "kubectl_manifest" "istio_operator_namespace" {
  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: istio-operator
  labels:
    istio-operator-managed: Reconcile
    istio-injection: disabled
YAML

  depends_on = [
    helm_release.istio_base
  ]
}

resource "helm_release" "istio_operator" {
  atomic           = true
  chart            = "./charts/istio-operator"
  create_namespace = true
  name             = "istio-operator"
  namespace        = kubectl_manifest.istio_operator_namespace.name
  wait_for_jobs    = true
}

resource "kubectl_manifest" "istio_config" {
  yaml_body = <<YAML
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: istio
  namespace: istio-system
spec:
  meshConfig:
    defaultConfig:
      proxyMetadata:
        # Enable basic DNS proxying
        ISTIO_META_DNS_CAPTURE: "true"
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: "${google_container_cluster.primary.name}"
      network: "${google_compute_network.vpc.name}"
  components:
    ingressGateways:
      - name: istio-eastwestgateway
        label:
          istio: eastwestgateway
          app: istio-eastwestgateway
        enabled: true
        k8s:
          env:
            # sni-dnat adds the clusters required for AUTO_PASSTHROUGH mode
            - name: ISTIO_META_ROUTER_MODE
              value: "sni-dnat"
          service:
            ports:
              - name: status-port
                port: 15021
                targetPort: 15021
              - name: tls
                port: 15443
                targetPort: 15443
              - name: tls-istiod
                port: 15012
                targetPort: 15012
              - name: tls-webhook
                port: 15017
                targetPort: 15017
YAML

  depends_on = [
    helm_release.istio_operator,
  ]
}

resource "kubectl_manifest" "eastwest_gateway" {
  yaml_body = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: istiod-gateway
spec:
  selector:
    istio: eastwestgateway
  servers:
    - port:
        name: tls-istiod
        number: 15012
        protocol: tls
      tls:
        mode: PASSTHROUGH        
      hosts:
        - "*"
    - port:
        name: tls-istiodwebhook
        number: 15017
        protocol: tls
      tls:
        mode: PASSTHROUGH          
      hosts:
        - "*"
YAML

  depends_on = [
    helm_release.istio_operator
  ]
}

resource "kubectl_manifest" "eastwest_gateway_vs" {
  yaml_body = <<YAML
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: istiod-vs
spec:
  hosts:
  - "*"
  gateways:
  - istiod-gateway
  tls:
  - match:
    - port: 15012
      sniHosts:
      - "*"
    route:
    - destination:
        host: istiod.istio-system.svc.cluster.local
        port:
          number: 15012
  - match:
    - port: 15017
      sniHosts:
      - "*"
    route:
    - destination:
        host: istiod.istio-system.svc.cluster.local
        port:
          number: 443
YAML

  depends_on = [
    helm_release.istio_operator
  ]
}

## Prometheus
resource "helm_release" "prometheus" {
  atomic           = true
  create_namespace = true
  name             = "prometheus"
  namespace        = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus"

  depends_on = [
    google_container_cluster.primary,
    google_container_node_pool.primary_nodes
  ]

  set {
    name  = "server.global.scrape_interval"
    value = "15s"
  }

  set {
    name  = "extraScrapeConfigs"
    value = local.scrape_configs
  }
}

## Kiali
resource "helm_release" "kiali" {
  atomic     = true
  name       = "kiali"
  namespace  = kubectl_manifest.istio_config.namespace
  repository = "https://kiali.org/helm-charts"
  chart      = "kiali-server"

  # This is for test purposes only. You should use a different auth strategy for production.
  set {
    name  = "auth.strategy"
    value = "anonymous"
  }

  set {
    name  = "external_services.prometheus.url"
    value = "http://prometheus-server.prometheus"
  }

  set {
    name  = "deployment.logger.log_level"
    value = "trace"
  }
}

resource "time_sleep" "wait_150_seconds" {
  create_duration = "150s"

  depends_on = [
    kubectl_manifest.eastwest_gateway
  ]
}

data "kubernetes_service" "eastwest-gateway-service" {
  metadata {
    name      = "istio-eastwestgateway"
    namespace = "istio-system"
  }

  # Give time for the service's load balancer to become active and get an external IP
  depends_on = [
    time_sleep.wait_150_seconds
  ]
}

data "kubernetes_service" "ingress-gateway-service" {
  metadata {
    name      = "istio-ingressgateway"
    namespace = "istio-system"
  }

  # Give time for the service's load balancer to become active and get an external IP
  depends_on = [
    time_sleep.wait_150_seconds
  ]
}
