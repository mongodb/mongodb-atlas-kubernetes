#!/bin/bash

set -e

image="$1"
timeout="$2"

echo "wait for pulling $image for $timeout min"

ok="OK"
command=$(docker pull "$image"  | awk -v RS="" '!/not found/{print "'"$ok"'"}' || true)

while [[ "$command" != "$ok" ]] && [[ "$time" -lt "$timeout" ]]; do
    echo "...wait for pulling $image"
    sleep 1m
    ((time = time + 3))
    command=$(docker pull "$image"  | awk -v RS="" '!/not found/{print "'"$ok"'"}' || true)
done

if [[ "$command" != "$ok" ]]; then
    exit 1
fi

echo OK
