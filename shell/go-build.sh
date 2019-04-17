#!/usr/bin/env bash

echo
echo "# Build Binary"
go build -mod=vendor -o /build/dt *.go