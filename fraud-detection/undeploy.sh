#!/usr/bin/env bash
set -e

kubectl delete -f fraud.yaml -n fraud-detection
kubectl delete -f bank.yaml -n fraud-detection
kubectl delete -f cards.yaml -n fraud-detection
kubectl delete -f accounts.yaml -n fraud-detection
kubectl delete -f insurance.yaml -n fraud-detection
kubectl delete -f policies.yaml -n fraud-detection
kubectl delete -f claims.yaml -n fraud-detection

kubectl delete namespace fraud-detection


