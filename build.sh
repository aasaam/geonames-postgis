#!/bin/bash
set -e

sudo apt install p7zip

docker rm -f geonames-postgis || true
docker rmi geonames-postgis || true

rm -rf var || true

# ./download.sh

docker pull postgis/postgis:13-3.1-alpine
docker run --name geonames-postgis --name geonames-postgis -p 5432:5432 -e POSTGRES_PASSWORD=geonames -e POSTGRES_USER=geonames -e POSTGRES_DB=geonames -d geonames-postgis

go build -o geodata
./geodata

docker stop geonames-postgis
docker cp geonames-postgis:/var/lib/postgresql/data var/data
docker build -f Dockerfile -t aasaam/geonames-postgis .
