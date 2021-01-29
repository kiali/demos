#!/usr/bin/env bash

if [[ "$#" -lt 1 || "$1" = "--help" ]]; then
	echo "Syntax: gen-services.sh <number of services>"
	echo ""
	exit
fi

NB_SERVICES="$1"
NAMESPACE="default"

LAST_ARG=""
for arg in "$@"
do
    if [[ "$LAST_ARG" = "-n" ]]; then
        NAMESPACE="$arg"
        LAST_ARG=""
    else
        LAST_ARG="$arg"
    fi
done

kubectl label namespace $NAMESPACE istio-injection=enabled

for (( c=0; c<$NB_SERVICES; c++ ))
do
    next=$(($c+1))
    if [[ $next -eq $NB_SERVICES ]]; then
        next=0
    fi
    cat "./deployment-tpl.yaml" \
        | sed -e "s:this-service:service-$c:g" \
        | sed -e "s:target-service:service-$next:g" \
        | sed -e "s:this-namespace:$NAMESPACE:g"
    echo "---"
done

