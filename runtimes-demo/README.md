# Runtimes Demo

🏃 🕔 📊

This demo can be used to start several services of various runtimes. Kiali will show custom metrics for these runtimes.

## Pre-requisite

- Kubernetes cluster running
- Better with Kiali and Istio installed

## Quick start

This quick start doesn't require you to clone the repo, but offers less interactivity.

If not already done, enable istio injection:

```bash
kubectl label namespace default istio-injection=enabled
```

Run as follow:

```bash
kubectl apply -f <(curl -L https://raw.githubusercontent.com/kiali/demos/master/runtimes-demo/quickstart.yml) -n default
```

## Advanced

- Clone this repo
- Read the fine manual!

```bash
make man
```

It covers a bunch of make targets, deployment options, etc.
