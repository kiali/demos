## Services Spawner

This is a placeholder for future demo.

Goal: have here some YAML templates + shell script that spawns N services (provided from command line).

For each service, service N talks to service N+1 (e.g. curl k8s probes).

That can perhaps be done without code container, just a busybox and some commands for probing.

The objective is to be able to easily spawn a high number of services that generate traffic.
