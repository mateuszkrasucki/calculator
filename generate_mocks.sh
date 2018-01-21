#!/bin/bash
set -e
set -u

for path in $(go list -f '{{.Dir}}' ./... | grep -v '/vendor/'); do
    for file in $(grep -le '^type [A-Z].* interface ' -- ${path}/*.go); do
        mock_path="$(dirname ${file})/mock_$(basename ${file%.*}).go"
        package="${PACKAGE:-"$(basename $(dirname "${file}"))"}"
        echo "Building mock ${mock_path}"
        mockgen -source="${file}" -package="${package}" >> $mock_path
    done
done
echo
