#!/bin/bash

set -e
set -x

if [ $# -eq 0 ]; then
    echo "Must supply an os and arch value"
    exit 1
fi

if [ -z "${1}" ]; then
    echo "Must supply an os value e.g. darwin"
    exit 1
fi

if [ -z "${2}" ]; then
    echo "Must supply an arch value e.g. amd64"
    exit 1
fi

os="${1}"
arch="${2}"

full_arch="${os}_${arch}"
binary="terraform-provider-gotemplate"
source="build/${full_arch}/${binary}"
output="${HOME}/.terraform.d/plugins/${full_arch}/"
mkdir -p "${output}"
cp "${source}" "${output}/${binary}"
