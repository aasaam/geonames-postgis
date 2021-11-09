<div align="center">
  <h1>
    GeoNames PostGIS
  </h1>
  <p>
    Prepossessed PostGIS database store GeoNames data
  </p>
  <p>
    <a href="https://github.com/aasaam/geonames-postgis/actions/workflows/ci.yml" target="_blank">
      <img src="https://github.com/aasaam/geonames-postgis/actions/workflows/ci.yml/badge.svg" alt="ci" />
    </a>
    <a href="https://hub.docker.com/r/aasaam/geonames-postgis" target="_blank"><img src="https://img.shields.io/docker/image-size/aasaam/geonames-postgis?label=docker%20image" alt="docker" /></a>
    <a href="https://github.com/aasaam/geonames-postgis/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/aasaam/geonames-postgis"></a>
  </p>
</div>

This database is pure latest stable version of postgres/postgis (14), that process data from [geonames.org](https://www.geonames.org/) and process countries, administrator code and geo names that featureClass is administrator area or population place.

Database will update schedule every week.

## Database structure

These are tables that data stored on it.

```sql
CREATE TABLE IF NOT EXISTS "countryInfo" (
  "geonameid" INT PRIMARY KEY,
  "continent" CHAR(2),
  "iso" CHAR(2) UNIQUE,
  "iso3" CHAR(3),
  "preferedLanguage" CHAR(2),
  "locales" VARCHAR(31)[],
  "tld" CHAR(3),
  "currency" CHAR(3),
  "area" BIGINT,
  "population" BIGINT,
  "neighbours" INT[],
  "polygons" geometry
);

CREATE TABLE IF NOT EXISTS "adminCode" (
  "id" INT PRIMARY KEY,
  "name" VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS "geo" (
  "geonameid" INT PRIMARY KEY,
  "name" VARCHAR(256),
  "country" INT,
  "adminCode" INT,
  "population" BIGINT,
  "timezone" VARCHAR(63),
  "location" geography(POINT, 4326),
  CONSTRAINT "geo_country__countryInfo_geonameid" FOREIGN KEY("country") REFERENCES "countryInfo"("geonameid"),
  CONSTRAINT "geo_adminCode__adminCode_id" FOREIGN KEY("adminCode") REFERENCES "adminCode"("id")
);
```

<div>
  <p align="center">
    <img alt="aasaam software development group" width="64" src="https://raw.githubusercontent.com/aasaam/information/master/logo/aasaam.svg">
    <br />
    aasaam software development group
  </p>
</div>
