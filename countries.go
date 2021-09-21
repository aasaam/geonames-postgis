package main

import (
	"database/sql"
	"log"

	pg "github.com/lib/pq"
)

type tmpCountryInfo struct {
	iso                string
	iso3               string
	isoNumeric         string
	fips               string
	country            string
	capital            string
	area               string
	population         string
	continent          string
	tld                string
	currencyCode       string
	currencyName       string
	phone              string
	postalcodeFormat   string
	postalCodeRegex    string
	languages          string
	geonameid          string
	neighbours         string
	equivalentFipsCode string
}

func processCountryInfo(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE "countryInfo" CASCADE;`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`
		SELECT
			"iso",
			"iso3",
			"isoNumeric",
			"fips",
			"country",
			"capital",
			"area",
			"population",
			"continent",
			"tld",
			"currencyCode",
			"currencyName",
			"phone",
			"postalcodeFormat",
			"postalCodeRegex",
			"languages",
			"geonameid",
			"neighbours",
			"equivalentFipsCode"
		FROM "tmp_countryInfo"
		WHERE "geonameid" != ''
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r tmpCountryInfo
		err := rows.Scan(
			&r.iso,
			&r.iso3,
			&r.isoNumeric,
			&r.fips,
			&r.country,
			&r.capital,
			&r.area,
			&r.population,
			&r.continent,
			&r.tld,
			&r.currencyCode,
			&r.currencyName,
			&r.phone,
			&r.postalcodeFormat,
			&r.postalCodeRegex,
			&r.languages,
			&r.geonameid,
			&r.neighbours,
			&r.equivalentFipsCode,
		)
		if err != nil {
			log.Fatal(err)
		}
		area := NewNullFloatInt(r.area)
		population := NewNullInt64(r.population)
		geonameid := toInt(r.geonameid)
		preferedLanguage := findPrefereLanguage(r.languages)
		locales := parseLocales(r.languages)
		polygons := getCountryPolygon(db, geonameid)
		neighbours := listOfNeighbours(r.neighbours)
		insertCountryInfo(
			db,
			r,
			geonameid,
			area,
			population,
			preferedLanguage,
			neighbours,
			locales,
			polygons,
		)

	}
}

func insertCountryInfo(
	db *sql.DB,
	tmpRow tmpCountryInfo,
	geonameid int,
	area sql.NullInt64,
	population sql.NullInt64,
	preferedLanguage string,
	neighbours []int,
	locales []string,
	polygons string,
) {
	var err error
	if len(polygons) > 1 {
		_, err = db.Exec(`
			INSERT INTO "countryInfo" (
				"geonameid",
				"neighbours",
				"currency",
				"locales",
				"preferedLanguage",
				"iso",
				"iso3",
				"polygons",
				"area",
				"population",
				"continent",
				"tld"
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				ST_GeomFromGeoJSON($8),
				$9,
				$10,
				$11,
				$12
			)
		`,
			geonameid,
			pg.Array(neighbours),
			NewNullString(tmpRow.currencyCode),
			pg.Array(locales),
			NewNullString(preferedLanguage),
			tmpRow.iso,
			tmpRow.iso3,
			polygons,
			area,
			population,
			tmpRow.continent,
			NewNullString(tmpRow.tld),
		)
	} else {
		_, err = db.Exec(`
		INSERT INTO "countryInfo" (
			"geonameid",
			"neighbours",
			"currency",
			"locales",
			"preferedLanguage",
			"iso",
			"iso3",
			"area",
			"population",
			"continent",
			"tld"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		)
	`,
			geonameid,
			pg.Array(neighbours),
			tmpRow.currencyCode,
			pg.Array(locales),
			preferedLanguage,
			tmpRow.iso,
			tmpRow.iso3,
			area,
			population,
			tmpRow.continent,
			tmpRow.tld,
		)
	}
	if err != nil {
		log.Fatal(err)
	}
}
