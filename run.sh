#!/bin/sh
docker run \
    -it \
    -v $(pwd):/go \
    -v /home/uki/Obrazy:/data \
    --rm \
    golang:1.8 \
    go run main.go
