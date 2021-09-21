#!/bin/bash

docker rm -f geonames-postgis-build-app || true
docker rm -f geonames-postgis || true
docker run --name geonames-postgis -p 5432:5432 -e POSTGRES_PASSWORD=geonames -e POSTGRES_USER=geonames -e POSTGRES_DB=geonames -v $(pwd)/init.db:/docker-entrypoint-initdb.d -v $(pwd)/var/geonames_extract:/geonames_extract -v $(pwd)/var/postgresql:/var/lib/postgresql/data -d postgis/postgis:13-3.1
docker build -f Dockerfile.build-app -t geonames-postgis-build-app .
docker run --name geonames-postgis-build-app --net=host geonames-postgis-build-app
