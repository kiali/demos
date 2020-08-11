#!/usr/bin/env bash
set -e

kubectl delete -f travel_control.yaml -n travel-control
kubectl delete -f travel_portal.yaml -n travel-portal
kubectl delete -f travel_agency.yaml -n travel-agency

kubectl delete namespace travel-agency
kubectl delete namespace travel-portal
kubectl delete namespace travel-control


