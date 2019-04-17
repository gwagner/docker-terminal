#!/bin/bash

echo
echo "# Building docker image"
docker build -f Dockerfile.ubuntu -t terminal:latest .

# Determine which home to use, inside or outside the container
DIR=$PWD
if [ ! -z $IS_CONTAINER ]; then
    DIR=$HOST_PROJECT_DIR
fi

echo
echo "# Building login shell"
docker run -it \
    -e "GOOS=darwin" \
    --mount type=bind,source="$DIR/shell",target="/workdir" -w "/workdir" \
    --mount type=bind,source="$DIR/bin",target="/build" -w "/workdir" \
    terminal:latest "go-build.sh"