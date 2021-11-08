#!/bin/bash

set -e

echo "Remove old database containers"
docker rm -f geonames-postgis || true
docker rmi geonames-postgis || true

echo "Bring data base up and store unclean data"
docker pull postgis/postgis:13-3.1-alpine
docker build --no-cache --tag geonames-postgis -f Dockerfile.db .
docker run --name geonames-postgis --name geonames-postgis \
  -v $(pwd)/init.db:/docker-entrypoint-initdb.d \
  -v $(pwd)/var/geonames_extract:/geonames_extract \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD=geonames -e POSTGRES_USER=geonames -e POSTGRES_DB=geonames \
  -d geonames-postgis

echo "Building process application"
go mod tidy
go build -o geodata
echo "Running process application"
./geodata

echo "Stop container"
docker stop geonames-postgis

echo "Copy container files"
docker cp geonames-postgis:/var/lib/postgresql/data var/data

echo "Make new image from stored data"
docker build -f Dockerfile -t aasaam/geonames-postgis .
