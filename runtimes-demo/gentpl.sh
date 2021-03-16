#!/usr/bin/env bash

if [[ "$#" -lt 1 || "$1" = "--help" ]]; then
	echo "Syntax: gentpl.sh <service name> <other options>"
	echo ""
	exit
fi

NAME="$1"
PULL_POLICY="Never"
DOMAIN=""
USER="jotak"
TAG="dev"
LAST_ARG=""

for arg in "$@"
do
    if [[ "$LAST_ARG" = "-pp" ]]; then
        PULL_POLICY="$arg"
        LAST_ARG=""
    elif [[ "$LAST_ARG" = "-d" ]]; then
        DOMAIN="$arg"
        LAST_ARG=""
    elif [[ "$LAST_ARG" = "-t" ]]; then
        TAG="$arg"
        LAST_ARG=""
    elif [[ "$LAST_ARG" = "-u" ]]; then
        USER="$arg"
        LAST_ARG=""
    else
        LAST_ARG="$arg"
    fi
done

IMAGE="${DOMAIN}${USER}/runtimes-$NAME:$TAG"
# Escape for sed
ESCAPED_IMAGE=$(printf '%s\n' "$IMAGE" | sed -e 's/[\/&]/\\&/g')


if [[ -f "./kube/$NAME.yml" ]] ; then
    cat "./kube/$NAME.yml" \
        | sed -e "s/#PULL_POLICY/$PULL_POLICY/g" \
        | sed -e "s/#IMAGE/$ESCAPED_IMAGE/g"
    echo "---"
fi
