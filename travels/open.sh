#!/usr/bin/env bash
set -e

# Linux helper
# xdg-open is a command to open a browser from the command line.
# Mac users can use "open" command
# Windows users can use "cmd" command
# Otherwise, the check is just to confirm you can access addons from your browser

xdg-open http://localhost:18080 &

kubectl port-forward svc/control 18080:8080 -n travel-control