#!/bin/bash

set -euo pipefail

max_retries=${MAX_RETRIES:-7}
backoff=${BACKOFF:-1}

retries=0
until (( retries == max_retries )) || "${@}"; do
    sleep "$(( (retries++)*backoff ))"
done
exit $?
