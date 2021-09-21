-- Copyright (c) 2021 aasaam software development group

-- enable postgis
CREATE EXTENSION postgis;

-- TEMPORARY TABLES

CREATE TABLE IF NOT EXISTS "tmp_geonameid" (
  "geonameid" INT,
  "name" VARCHAR(256),
  "asciiname" VARCHAR(256),
  "alternatenames" TEXT,
  "latitude" VARCHAR(256),
  "longitude" VARCHAR(256),
  "featureClass" VARCHAR(4),
  "featureCode" VARCHAR(16),
  "countryCode" VARCHAR(2),
  "cc2" VARCHAR(256),
  "admin1Code" VARCHAR(24),
  "admin2Code" VARCHAR(96),
  "admin3Code" VARCHAR(24),
  "admin4Code" VARCHAR(24),
  "population" VARCHAR(64),
  "elevation" VARCHAR(64),
  "dem" VARCHAR(64),
  "timezone" VARCHAR(64),
  "modificationDate" DATE
);

CREATE INDEX concurrently "tmp_geonameid_geonameid"
ON "tmp_geonameid" USING btree ("geonameid");

CREATE INDEX concurrently "tmp_geonameid_countryCode"
ON "tmp_geonameid" USING btree ("countryCode");

CREATE TABLE IF NOT EXISTS "tmp_countryInfo" (
  "iso" TEXT,
  "iso3" TEXT,
  "isoNumeric" TEXT,
  "fips" TEXT,
  "country" TEXT,
  "capital" TEXT,
  "area" TEXT,
  "population" TEXT,
  "continent" TEXT,
  "tld" TEXT,
  "currencyCode" TEXT,
  "currencyName" TEXT,
  "phone" TEXT,
  "postalcodeFormat" TEXT,
  "postalCodeRegex" TEXT,
  "languages" TEXT,
  "geonameid" TEXT,
  "neighbours" TEXT,
  "equivalentFipsCode" TEXT
);

CREATE TABLE IF NOT EXISTS "tmp_shapesAllLow" (
  "geonameid_1" INT,
  "polygonData" TEXT
);

CREATE TABLE IF NOT EXISTS "tmp_ready" (
  "field1" INT
);

-- PRODUCTION DATABASE

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

CREATE INDEX "countryInfo_polygons"
ON "countryInfo" USING GIST ("polygons");

CREATE TABLE IF NOT EXISTS "geo" (
  "geonameid" INT PRIMARY KEY,
  "name" VARCHAR(256),
  "featureClass" VARCHAR(5),
  "featureCode" VARCHAR(15),
  "country" INT,
  "admin1Code" VARCHAR(23),
  "admin2Code" VARCHAR(95),
  "admin3Code" VARCHAR(23),
  "admin4Code" VARCHAR(23),
  "population" BIGINT,
  "elevation" INT,
  "dem" INT,
  "timezone" VARCHAR(63),
  "location" geography(POINT, 4326),
  CONSTRAINT "geo_country__countryInfo_geonameid" FOREIGN KEY("country") REFERENCES "countryInfo"("geonameid")
);

CREATE INDEX concurrently "geo_featureClass"
ON "geo" USING btree ("featureClass");

CREATE INDEX concurrently "geo_featureCode"
ON "geo" USING btree ("featureCode");

CREATE INDEX concurrently "geo_admin1Code"
ON "geo" USING btree ("admin1Code");

CREATE INDEX "geo_location"
ON "geo" USING GIST ("location");

-- IMPORT TEMPORARY DATA
COPY "tmp_shapesAllLow" FROM '/geonames_extract/shapes_all_low.tsv' DELIMITER E'\t';
COPY "tmp_countryInfo" FROM '/geonames_extract/countryInfo.tsv' DELIMITER E'\t';
COPY "tmp_geonameid" FROM '/geonames_extract/geonameid.tsv' DELIMITER E'\t';

INSERT INTO "tmp_ready" ("field1") VALUES (1);
