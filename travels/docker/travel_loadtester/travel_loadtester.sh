#!/usr/bin/env bash

CONNECTIONS=${CONNECTIONS:-10}
NUMCALLS=${NUMCALLS:-20}
RATE=${RATE:-0}
SLEEP=${SLEEP:-5}

echo "Travel Portal LoadTester Configuration"
echo "================================================"
echo "CONNECTIONS=${CONNECTIONS}"
echo "NUMCALLS=${NUMCALLS}"
echo "RATE=${RATE}"
echo "SLEEP=${SLEEP}"
echo "TRAVELS_AGENCY_SERVICE=${TRAVELS_AGENCY_SERVICE}"
echo "================================================"

if [ -z "${TRAVELS_AGENCY_SERVICE}" ]
then
    echo "TRAVELS_AGENCY_SERVICE must be non empty"
    exit 1
fi

while true
do
    NOW=$(date +"%Y/%m/%d_%H:%M:%S")
    echo "[${NOW}][Sleep ${SLEEP}] /usr/bin/fortio load -c ${CONNECTIONS} -qps ${RATE} -n ${NUMCALLS} -loglevel Warning ${TRAVELS_AGENCY_SERVICE}"
    /usr/bin/fortio load -c ${CONNECTIONS} -qps ${RATE} -n ${NUMCALLS} -loglevel Warning "${TRAVELS_AGENCY_SERVICE}/travels"
    sleep ${SLEEP}
done