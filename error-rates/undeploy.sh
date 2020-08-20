#!/usr/bin/env bash
set -e

kubectl delete -f alpha.yaml -n alpha
kubectl delete -f beta.yaml -n beta

kubectl delete namespace alpha
kubectl delete namespace beta


