#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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
