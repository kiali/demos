#!/bin/bash

# Set up your host machine for CodeReady Containers
crc setup

# Start the CodeReady Containers virtual machine
crc start --memory 16384

# Add the cached oc executable to your PATH:
eval $(crc oc-env)

# Log in as a cluster admin user to install the operators needed
oc config use-context crc-admin

# Subscribe to the required operators
oc apply -f ossm-subs.yaml

# Install two independent meshes (east and west)
oc create namespace east-mesh-system
oc apply -n east-mesh-system -f east/east-ossm.yaml
oc create namespace west-mesh-system
oc apply -n west-mesh-system -f west/west-ossm.yaml

# Create the namespaces for the Travels application in both meshes:
oc create namespace east-travel-agency
oc create namespace east-travel-portal
oc create namespace east-travel-control
oc create namespace west-travel-agency

# Wait for both control planes to be ready
oc wait --for condition=Ready -n east-mesh-system smmr/default --timeout 300s
oc wait --for condition=Ready -n west-mesh-system smmr/default --timeout 300s

# Create in each mesh namespace a configmap containing a root certificate that is used to validate client certificates in the trust domain used by the other mesh
oc get configmap istio-ca-root-cert -o jsonpath='{.data.root-cert\.pem}' -n east-mesh-system > east-cert.pem
oc create configmap east-ca-root-cert --from-file=root-cert.pem=east-cert.pem -n west-mesh-system
oc get configmap istio-ca-root-cert -o jsonpath='{.data.root-cert\.pem}' -n west-mesh-system > west-cert.pem
oc create configmap west-ca-root-cert --from-file=root-cert.pem=west-cert.pem -n east-mesh-system

# Configure all the federation resources for travels application
oc apply -n east-mesh-system -f east/east-federation.yaml
oc apply -n west-mesh-system -f west/west-federation.yaml

# Create the application resources
oc apply -n east-travel-agency -f east/east-travel-agency.yaml
oc apply -n east-travel-portal -f east/east-travel-portal.yaml
oc apply -n east-travel-control -f east/east-travel-control.yaml
oc apply -n west-travel-agency -f west/west-travel-agency.yaml

echo "Installation complete.

Visit both Kiali instances to look for the complete topology:

Links:
East Kiali: https://kiali-east-mesh-system.apps-crc.testing
West Kiali: https://kiali-west-mesh-system.apps-crc.testing

Log in with kubeadmin user:

$(crc console --credentials)
"

