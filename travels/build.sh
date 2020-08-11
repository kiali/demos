#!/usr/bin/env bash
set -e

DOCKER_VERSION=v2

## Travel Control

DOCKER_TRAVEL_CONTROL=lucasponce/travel_control
DOCKER_TRAVEL_CONTROL_TAG=${DOCKER_TRAVEL_CONTROL}:${DOCKER_VERSION}

rm -Rf docker/travel_control/travel_control docker/travel_control/html
cd travel_control
go build -o ../docker/travel_control/travel_control
cp -R html ../docker/travel_control
cd ..

podman build -t ${DOCKER_TRAVEL_CONTROL_TAG} docker/travel_control

## Travel Portal

DOCKER_TRAVEL_PORTAL=lucasponce/travel_portal
DOCKER_TRAVEL_PORTAL_TAG=${DOCKER_TRAVEL_PORTAL}:${DOCKER_VERSION}

rm -Rf docker/travel_portal/travel_portal
cd travel_portal
go build -o ../docker/travel_portal/travel_portal
cd ..

podman build -t ${DOCKER_TRAVEL_PORTAL_TAG} docker/travel_portal

## MySQL

DOCKER_TRAVEL_MYSQL=lucasponce/mysqldb
DOCKER_TRAVEL_MYSQL_TAG=${DOCKER_TRAVEL_MYSQL}:${DOCKER_VERSION}

podman build -t ${DOCKER_TRAVEL_MYSQL_TAG} docker/travel_agency/mysql

## Cars

DOCKER_TRAVEL_CARS=lucasponce/cars
DOCKER_TRAVEL_CARS_TAG=${DOCKER_TRAVEL_CARS}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/cars/cars
cd travel_agency/cars
go build -o ../../docker/travel_agency/cars/cars
cd ../..

podman build -t ${DOCKER_TRAVEL_CARS_TAG} docker/travel_agency/cars

## Discounts

DOCKER_TRAVEL_DISCOUNTS=lucasponce/discounts
DOCKER_TRAVEL_DISCOUNTS_TAG=${DOCKER_TRAVEL_DISCOUNTS}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/discounts/discounts
cd travel_agency/discounts
go build -o ../../docker/travel_agency/discounts/discounts
cd ../..

podman build -t ${DOCKER_TRAVEL_DISCOUNTS_TAG} docker/travel_agency/discounts

## Flights

DOCKER_TRAVEL_FLIGHTS=lucasponce/flights
DOCKER_TRAVEL_FLIGHTS_TAG=${DOCKER_TRAVEL_FLIGHTS}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/flights/flights
cd travel_agency/flights
go build -o ../../docker/travel_agency/flights/flights
cd ../..

podman build -t ${DOCKER_TRAVEL_FLIGHTS_TAG} docker/travel_agency/flights

## Hotels

DOCKER_TRAVEL_HOTELS=lucasponce/hotels
DOCKER_TRAVEL_HOTELS_TAG=${DOCKER_TRAVEL_HOTELS}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/hotels/hotels
cd travel_agency/hotels
go build -o ../../docker/travel_agency/hotels/hotels
cd ../..

podman build -t ${DOCKER_TRAVEL_HOTELS_TAG} docker/travel_agency/hotels

## Insurances

DOCKER_TRAVEL_INSURANCES=lucasponce/insurances
DOCKER_TRAVEL_INSURANCES_TAG=${DOCKER_TRAVEL_INSURANCES}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/insurances/insurances
cd travel_agency/insurances
go build -o ../../docker/travel_agency/insurances/insurances
cd ../..

podman build -t ${DOCKER_TRAVEL_INSURANCES_TAG} docker/travel_agency/insurances

## Travels

DOCKER_TRAVEL_TRAVELS=lucasponce/travels
DOCKER_TRAVEL_TRAVELS_TAG=${DOCKER_TRAVEL_TRAVELS}:${DOCKER_VERSION}

rm -Rf docker/travel_agency/travels/travels
cd travel_agency/travels
go build -o ../../docker/travel_agency/travels/travels
cd ../..

podman build -t ${DOCKER_TRAVEL_TRAVELS_TAG} docker/travel_agency/travels

podman login docker.io
podman push ${DOCKER_TRAVEL_CONTROL_TAG}
podman push ${DOCKER_TRAVEL_PORTAL_TAG}
podman push ${DOCKER_TRAVEL_MYSQL_TAG}
podman push ${DOCKER_TRAVEL_CARS_TAG}
podman push ${DOCKER_TRAVEL_DISCOUNTS_TAG}
podman push ${DOCKER_TRAVEL_FLIGHTS_TAG}
podman push ${DOCKER_TRAVEL_HOTELS_TAG}
podman push ${DOCKER_TRAVEL_INSURANCES_TAG}
podman push ${DOCKER_TRAVEL_TRAVELS_TAG}