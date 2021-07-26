locals {
  workloads = {
    productpage-v1 = {
      app     = "productpage",
      version = "v1"
    }
    details-v1 = {
      app     = "details",
      version = "v1"
    }
    ratings-v1 = {
      app     = "ratings",
      version = "v1"
    }
    reviews-v1 = {
      app     = "reviews",
      version = "v1"
    }
    reviews-v2 = {
      app     = "reviews",
      version = "v2"
    }
    reviews-v3 = {
      app     = "reviews",
      version = "v3"
    }
  }
  # Just using the host network for the container. It's the most straightforward option.
  # It's also possible to use the default bridge network and bind to localhost but that
  # requires a Sidecar config with the defaultEndpoint set to localhost to tell the 
  # workload's envoy proxy to send traffic there instead of '0.0.0.0'. There also seems
  # to be some weird behavior between the docker iptables rules for the bridge networking
  # and the envoy ones.
  workload-services = { for workload_key, workload in local.workloads : workload_key => <<EOF
[Unit]
Description=${title(workload.app)}
After=docker.service
Requires=docker.service

[Service]
ExecStartPre=-/usr/bin/docker stop ${workload.app}
ExecStartPre=-/usr/bin/docker rm ${workload.app}
ExecStartPre=/usr/bin/docker pull istio/examples-bookinfo-${workload.app}-${workload.version}:latest
ExecStart=/usr/bin/docker run --net host --rm --name ${workload.app} istio/examples-bookinfo-${workload.app}-${workload.version}:latest
TimeoutStartSec=0
Restart=always

[Install]
WantedBy=multi-user.target

EOF
  }

  # All the values are base64 because terraform string interpolation was doing
  # weird things to the yaml formatted strings. There are much better ways of
  # passing these values for a real application.
  #
  # Starting istio last is important since there were some problems with 
  # cloud dns when istio was started before docker.
  workload_startup_scripts = { for workload_key in keys(local.workloads) : workload_key => <<EOF
#!/bin/bash

set -x

curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io

printf "%s" "${base64encode(local.workload-services[workload_key])}" | base64 --decode > /etc/systemd/system/docker.${workload_key}.service
systemctl start docker.${workload_key}
  
curl -LO https://storage.googleapis.com/istio-release/releases/1.10.1/deb/istio-sidecar.deb
dpkg -i istio-sidecar.deb
mkdir -p /etc/istio/proxy /etc/certs /var/run/secrets/tokens
printf "%s" "${data.local_file.root_cert[workload_key].content_base64}" | base64 --decode > /etc/certs/root-cert.pem 
printf "%s" "${data.local_file.istio_token[workload_key].content_base64}" | base64 --decode > /var/run/secrets/tokens/istio-token 
printf "%s" "${data.local_file.cluster_env[workload_key].content_base64}" | base64 --decode > /var/lib/istio/envoy/cluster.env 
printf "%s" "${data.local_file.mesh[workload_key].content_base64}" | base64 --decode > /etc/istio/config/mesh
chown -R istio-proxy /var/lib/istio /etc/certs /etc/istio/proxy /etc/istio/config /var/run/secrets /etc/certs/root-cert.pem
systemctl start istio

EOF
  }

  workloads_with_vm_info = {
    for vm in google_compute_instance.workload_vm : vm.name => merge(vm, local.workloads[vm.name])
  }

  apps = toset([for workload in values(local.workloads) : workload.app])

  common_vm_tags = ["app"]

  ssh_key = var.ssh_key == "" ? file(var.ssh_key_filepath) : var.ssh_key

  # The workload names can be used here as scrape targets because
  # Google Cloud DNS creates internal hostnames for each of the vms.
  scrape_configs = <<YAML
- job_name: bookinfo-vms
  metrics_path: '/stats/prometheus'
  static_configs:
  - targets:
%{for workload in keys(local.workloads)~}
    - '${workload}:15020'
%{endfor}
YAML
}