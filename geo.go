package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type tmp_geonameid struct {
	geonameid        int
	name             string
	asciiname        string
	latitude         string
	longitude        string
	featureClass     string
	featureCode      string
	countryCode      string
	admin1Code       string
	admin2Code       string
	admin3Code       string
	admin4Code       string
	population       string
	elevation        string
	dem              string
	timezone         string
	modificationDate string
}

type GeoNameLang map[string]string

func getProgress(db *sql.DB) int {
	rows, err := db.Query(`
		SELECT
			COUNT(*)
		FROM "tmp_geonameid"
		WHERE "featureClass" = 'A' OR "featureClass" = 'P'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var num int = 0
	for rows.Next() {
		err2 := rows.Scan(&num)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	return num
}

func buildGeoData(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE "geo" CASCADE;`)
	if err != nil {
		log.Fatal(err)
	}

	totalNumbers := getProgress(db)

	rows, err := db.Query(`
		SELECT
			"geonameid",
			"name",
			"asciiname",
			"latitude",
			"longitude",
			"featureClass",
			"featureCode",
			"countryCode",
			"admin1Code",
			"admin2Code",
			"admin3Code",
			"admin4Code",
			"population",
			"elevation",
			"dem",
			"timezone",
			"modificationDate"
		FROM "tmp_geonameid"
		WHERE "featureClass" = 'A' OR "featureClass" = 'P'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		i++
		var r tmp_geonameid

		err2 := rows.Scan(
			&r.geonameid,
			&r.name,
			&r.asciiname,
			&r.latitude,
			&r.longitude,
			&r.featureClass,
			&r.featureCode,
			&r.countryCode,
			&r.admin1Code,
			&r.admin2Code,
			&r.admin3Code,
			&r.admin4Code,
			&r.population,
			&r.elevation,
			&r.dem,
			&r.timezone,
			&r.modificationDate,
		)

		if (i % 1000) == 0 {
			fmt.Printf("Current %d/%d\n", i, totalNumbers)
		}

		if err2 == nil {

			countryGeoId, isCountryCodeExist := isoToGeonameID[r.countryCode]
			if !isCountryCodeExist {
				continue
			}

			adminCode1Name := r.countryCode + "." + r.admin1Code
			adminCode1ID, ok := adminCodeMapList[adminCode1Name]

			if !ok {
				adminCode1ID = 0
			}

			population := NewNullInt64(r.population)

			elevation := NewNullInt64(r.elevation)
			dem := NewNullInt64(r.dem)

			processGeoFields(db, &r, adminCode1ID, countryGeoId, elevation, population, dem)

		}
	}
}

func processGeoFields(
	db *sql.DB,
	tmpRow *tmp_geonameid,
	adminCode1ID int,
	countryGeoId int,
	elevation sql.NullInt64,
	population sql.NullInt64,
	dem sql.NullInt64,
) {
	_, err := db.Exec(`
		INSERT INTO "geo" (
			"geonameid",
			"location",
			"name",
			"country",
			"adminCode",
			"population",
			"timezone"
		) VALUES(
			$1,
			ST_GeomFromText($2),
			$3,
			$4,
			$5,
			$6,
			$7
		)
	`,
		tmpRow.geonameid,
		fmt.Sprintf("Point(%s %s)", tmpRow.longitude, tmpRow.latitude),
		fixName(tmpRow.asciiname),
		countryGeoId,
		NewNullInt64FromInt(adminCode1ID),
		population,
		tmpRow.timezone,
	)
	if err != nil {
		log.Fatal(err)
	}
}
