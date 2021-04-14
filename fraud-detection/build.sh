#!/usr/bin/env bash
set -e

DORP=${DORP:-podman}

DOCKER_VERSION=v1

IMAGE_HUB=quay.io

## Accounts

DOCKER_FRAUD_ACCOUNTS=${IMAGE_HUB}/kiali/demo_fraud_accounts
DOCKER_FRAUD_ACCOUNTS_TAG=${DOCKER_FRAUD_ACCOUNTS}:${DOCKER_VERSION}

rm -Rf docker/accounts/accounts
cd accounts_server
go build -o ../docker/accounts/accounts
cd ..

${DORP} build -t ${DOCKER_FRAUD_ACCOUNTS_TAG} docker/accounts

## Bank

DOCKER_FRAUD_BANK=${IMAGE_HUB}/kiali/demo_fraud_bank
DOCKER_FRAUD_BANK_TAG=${DOCKER_FRAUD_BANK}:${DOCKER_VERSION}

rm -Rf docker/bank/bank
cd bank_server
go build -o ../docker/bank/bank
cd ..

${DORP} build -t ${DOCKER_FRAUD_BANK_TAG} docker/bank

## Cards

DOCKER_FRAUD_CARDS=${IMAGE_HUB}/kiali/demo_fraud_cards
DOCKER_FRAUD_CARDS_TAG=${DOCKER_FRAUD_CARDS}:${DOCKER_VERSION}

rm -Rf docker/cards/cards
cd cards_server
go build -o ../docker/cards/cards
cd ..

${DORP} build -t ${DOCKER_FRAUD_CARDS_TAG} docker/cards

## Claims

DOCKER_FRAUD_CLAIMS=${IMAGE_HUB}/kiali/demo_fraud_claims
DOCKER_FRAUD_CLAIMS_TAG=${DOCKER_FRAUD_CLAIMS}:${DOCKER_VERSION}

rm -Rf docker/claims/claims
cd claims_server
go build -o ../docker/claims/claims
cd ..

${DORP} build -t ${DOCKER_FRAUD_CLAIMS_TAG} docker/claims

## Fraud

DOCKER_FRAUD_FRAUD=${IMAGE_HUB}/kiali/demo_fraud_fraud
DOCKER_FRAUD_FRAUD_TAG=${DOCKER_FRAUD_FRAUD}:${DOCKER_VERSION}

rm -Rf docker/fraud/fraud
cd fraud
go build -o ../docker/fraud/fraud
cd ..

${DORP} build -t ${DOCKER_FRAUD_FRAUD_TAG} docker/fraud

## Insurance

DOCKER_FRAUD_INSURANCE=${IMAGE_HUB}/kiali/demo_fraud_insurance
DOCKER_FRAUD_INSURANCE_TAG=${DOCKER_FRAUD_INSURANCE}:${DOCKER_VERSION}

rm -Rf docker/insurance/insurance
cd insurance_server
go build -o ../docker/insurance/insurance
cd ..

${DORP} build -t ${DOCKER_FRAUD_INSURANCE_TAG} docker/insurance

## Policies

DOCKER_FRAUD_POLICIES=${IMAGE_HUB}/kiali/demo_fraud_policies
DOCKER_FRAUD_POLICIES_TAG=${DOCKER_FRAUD_POLICIES}:${DOCKER_VERSION}

rm -Rf docker/policies/policies
cd policies_server
go build -o ../docker/policies/policies
cd ..

${DORP} build -t ${DOCKER_FRAUD_POLICIES_TAG} docker/policies


## Push images

${DORP} login ${IMAGE_HUB}
${DORP} push ${DOCKER_FRAUD_ACCOUNTS_TAG}
${DORP} push ${DOCKER_FRAUD_BANK_TAG}
${DORP} push ${DOCKER_FRAUD_CARDS_TAG}
${DORP} push ${DOCKER_FRAUD_CLAIMS_TAG}
${DORP} push ${DOCKER_FRAUD_FRAUD_TAG}
${DORP} push ${DOCKER_FRAUD_INSURANCE_TAG}
${DORP} push ${DOCKER_FRAUD_POLICIES_TAG}
