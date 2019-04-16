#!/bin/bash

echo
echo "# Building docker image"
docker build -f Dockerfile.ubuntu -t terminal:latest .

echo
echo "# Building login shell"
docker run -it \
    -e "GOOS=darwin" \
    --mount type=bind,source="$PWD/shell",target="/workdir" -w "/workdir" \
    --mount type=bind,source="$PWD/bin",target="/build" -w "/workdir" \
    terminal:latest "go-build.sh"