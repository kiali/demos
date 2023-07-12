#!/usr/bin/env bash
set -e

kubectl create namespace alpha
kubectl create namespace beta
kubectl create namespace gamma

kubectl label namespace alpha istio-injection=enabled
kubectl label namespace beta istio-injection=enabled
kubectl label namespace gamma istio-injection=enabled

kubectl apply -f alpha.yaml -n alpha
kubectl apply -f beta.yaml -n beta
kubectl apply -f beta.yaml -n gamma
