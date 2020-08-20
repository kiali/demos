#!/usr/bin/env bash
set -e

DORP=${DORP:-podman}

DOCKER_VERSION=v1

## Server

DOCKER_SERVER=kiali/demo_error_rates_server
DOCKER_SERVER_TAG=${DOCKER_SERVER}:${DOCKER_VERSION}

rm -Rf docker/server/server
cd server
go build -o ../docker/server/server
cd ..

${DORP} build -t ${DOCKER_SERVER_TAG} docker/server

## Client

DOCKER_CLIENT=kiali/demo_error_rates_client
DOCKER_CLIENT_TAG=${DOCKER_CLIENT}:${DOCKER_VERSION}

rm -Rf docker/client/client
cd client
go build -o ../docker/client/client
cd ..

${DORP} build -t ${DOCKER_CLIENT_TAG} docker/client


${DORP} login docker.io
${DORP} push ${DOCKER_SERVER_TAG}
${DORP} push ${DOCKER_CLIENT_TAG}
