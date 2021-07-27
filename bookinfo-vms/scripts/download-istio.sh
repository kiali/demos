#!/bin/bash

# This script downloads the latest version of istio for the purpose of
# having the helm charts in the bookinfo-vms/charts dir. It is called
# by terraform but it can also be used on its own.

set -e

CHARTS_DIR=$1
if [ -z "$CHARTS_DIR" ]; then
    CHARTS_DIR=$0
fi

curl -L https://istio.io/downloadIstio | sh - &>/dev/null

rm -rf $CHARTS_DIR/charts
ISTIO_DIR=$(ls -d */ | grep "istio-")
mv $ISTIO_DIR/manifests/charts/ $CHARTS_DIR/charts
# Remove the namespace chart since this is managed in terraform.
rm $CHARTS_DIR/charts/istio-operator/templates/namespace.yaml
rm -rf $ISTIO_DIR
