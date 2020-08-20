#!/usr/bin/env bash
set -e

kubectl create namespace alpha
kubectl create namespace betha

kubectl label namespace alpha istio-injection=enabled
kubectl label namespace betha istio-injection=enabled

kubectl apply -f alpha.yaml -n alpha
kubectl apply -f betha.yaml -n betha
