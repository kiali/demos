#!/usr/bin/env bash
set -e

kubectl create namespace fraud-detection

kubectl label namespace fraud-detection istio-injection=enabled

kubectl apply -f accounts.yaml -n fraud-detection
kubectl apply -f cards.yaml -n fraud-detection
kubectl apply -f bank.yaml -n fraud-detection
kubectl apply -f policies.yaml -n fraud-detection
kubectl apply -f claims.yaml -n fraud-detection
kubectl apply -f insurance.yaml -n fraud-detection
kubectl apply -f fraud.yaml -n fraud-detection
