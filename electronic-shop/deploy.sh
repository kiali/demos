#!/usr/bin/env bash
set -e

kubectl create namespace electronic-shop
kubectl label namespace electronic-shop istio-injection=enabled
kubectl apply -f electronic-shop.yaml -n electronic-shop

