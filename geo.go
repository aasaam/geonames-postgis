package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
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

func getProgress(db *sql.DB) *pb.ProgressBar {
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

	// create bar
	bar := pb.New(num)

	// refresh info every second (default 200ms)
	bar.SetRefreshRate(time.Second * 5)
	bar.SetWriter(os.Stdout)

	bar.Start()

	return bar
}

func buildGeoData(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE "geo" CASCADE;`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ready to insert...")
	progressBar := getProgress(db)

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

	for rows.Next() {
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
		progressBar.Increment()
		if err2 == nil {

			countryGeoId, isCountryCodeExist := isoToGeonameID[r.countryCode]
			if !isCountryCodeExist {
				continue
			}

			population := NewNullInt64(r.population)

			elevation := NewNullInt64(r.elevation)
			dem := NewNullInt64(r.dem)

			processGeoFields(db, &r, countryGeoId, elevation, population, dem)

		}
	}
}

func processGeoFields(
	db *sql.DB,
	tmpRow *tmp_geonameid,
	countryGeoId int,
	elevation sql.NullInt64,
	population sql.NullInt64,
	dem sql.NullInt64,
) {
	_, err := db.Exec(`
		INSERT INTO "geo" (
			"geonameid",
			"name",
			"location",
			"featureClass",
			"featureCode",
			"country",
			"admin1Code",
			"admin2Code",
			"admin3Code",
			"admin4Code",
			"population",
			"elevation",
			"dem",
			"timezone"
		) VALUES(
			$1,
			$2,
			ST_GeomFromText($3),
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14
		)
	`,
		tmpRow.geonameid,
		fixName(tmpRow.asciiname),
		fmt.Sprintf("Point(%s %s)", tmpRow.longitude, tmpRow.latitude),
		tmpRow.featureClass,
		tmpRow.featureCode,
		countryGeoId,
		tmpRow.admin1Code,
		tmpRow.admin2Code,
		tmpRow.admin3Code,
		tmpRow.admin4Code,
		population,
		elevation,
		dem,
		tmpRow.timezone,
	)
	if err != nil {
		log.Fatal(err)
	}
}
