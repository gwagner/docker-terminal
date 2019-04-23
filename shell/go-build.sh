#!/usr/bin/env bash

if [ ! -d "vendor/" ]; then
    echo
    echo "# Vendor Modules"
    go mod vendor
fi

echo
echo "# Build Binary"
go build -mod=vendor -o /build/dt *.go
