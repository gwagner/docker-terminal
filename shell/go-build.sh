#!/usr/bin/env bash

echo
echo "# Vendor Modules"
go mod vendor

echo
echo "# Build Binary"
go build -mod=vendor -o /build/dt *.go
