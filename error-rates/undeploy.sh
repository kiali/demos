#!/usr/bin/env bash
set -e

kubectl delete -f alpha.yaml -n alpha
kubectl delete -f betha.yaml -n betha

kubectl delete namespace alpha
kubectl delete namespace betha


