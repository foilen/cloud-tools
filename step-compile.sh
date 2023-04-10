#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Compile ]----
go build -o build/bin/az-dns-update ./az-dns-update
chmod +x build/bin/az-dns-update
