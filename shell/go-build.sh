#!/usr/bin/env bash

echo "# Install Dependencies"
echo
go mod vendor

echo "# Build Binary"
echo
go build -o /build/dt *.go