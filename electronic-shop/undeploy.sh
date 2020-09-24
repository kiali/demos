#!/usr/bin/env bash
set -e

kubectl delete -f electronic-shop.yaml -n electronic-shop
kubectl delete namespace electronic-shop


