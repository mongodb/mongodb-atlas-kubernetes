#!/bin/bash

set -euo pipefail

dbpath=${DB_PATH:-/tmp/mdb}
dblogs=${DB_LOGS:-/tmp/mdb.log}
import_dir=${IMPORT_DIR:-imports}
collections=${COLLECTIONS:-flaky}

mkdir -p "${import_dir}"
mongod --dbpath="${dbpath}" --logpath="${dblogs}" --fork

for collection in ${collections}; do
	mongoimport --collection="${collection}" \
		--file "${import_dir}/${collection}.json" --jsonArray
done

mongosh

