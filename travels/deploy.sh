#!/usr/bin/env bash
set -e

kubectl create namespace travel-agency
kubectl create namespace travel-portal
kubectl create namespace travel-control

kubectl label namespace travel-agency istio-injection=enabled
kubectl label namespace travel-portal istio-injection=enabled
kubectl label namespace travel-control istio-injection=enabled

kubectl apply -f travel_agency.yaml -n travel-agency
kubectl apply -f travel_portal.yaml -n travel-portal
kubectl apply -f travel_control.yaml -n travel-control
