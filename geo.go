package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/cheggaaa/pb/v3"
	_ "github.com/lib/pq"
	"golang.org/x/text/language"
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

type tmp_hierarchy struct {
	geonameid_1 int
	geonameid_2 int
	relation    string
}

type tmp_alternateNamesV2 struct {
	geonameid   int
	isolanguage string
	name        string
}

type GeoNameLang map[string]string

var faRegexMatch = regexp.MustCompile(`[آابپتثجچحخدذرزژسشصضطظعغفقکگلمنوهیيك]+`)
var faRegexReplace = regexp.MustCompile(`[^آابپتثجچحخدذرزژسشصضطظعغفقکگلمنوهیيك ]`)

func findAlternativeNames(db *sql.DB, geonameid int) map[string]string {
	translate := make(map[string]string)
	rows, err := db.Query(`
		SELECT
			"geonameid",
			"isolanguage",
			"name"
		FROM "tmp_alternateNamesV2" WHERE "geonameid" = $1 ORDER BY "alternateNameId" ASC, "isPreferredName" ASC
	`, geonameid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// langData := make(map[string]map[string]float64)

	for rows.Next() {
		var r tmp_alternateNamesV2
		rows.Scan(
			&r.geonameid,
			&r.isolanguage,
			&r.name,
		)

		if _, ok := translate[r.isolanguage]; ok {
			continue
		}

		tag, err3 := language.ParseBase(r.isolanguage)
		if err3 == nil && tag.String() == r.isolanguage {
			if r.isolanguage == "fa" {
				if faRegexMatch.MatchString(r.name) {
					translate[r.isolanguage] = faRegexReplace.ReplaceAllString(r.name, "")
				}
			} else {
				translate[r.isolanguage] = r.name
			}
		}

	}

	return translate
}

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

		if err2 == nil {

			countryGeoId, isCountryCodeExist := isoToGeonameID[r.countryCode]
			if !isCountryCodeExist {
				continue
			}

			translates := findAlternativeNames(db, r.geonameid)
			names, _ := json.Marshal(translates)

			population := NewNullInt64(r.population)

			elevation := NewNullInt64(r.elevation)
			dem := NewNullInt64(r.dem)

			processGeoFields(db, &r, countryGeoId, names, elevation, population, dem)
			progressBar.Increment()
		}
	}
}

func processGeoFields(
	db *sql.DB,
	tmpRow *tmp_geonameid,
	countryGeoId int,
	names []byte,
	elevation sql.NullInt64,
	population sql.NullInt64,
	dem sql.NullInt64,
) {
	_, err := db.Exec(`
		INSERT INTO "geo" (
			"geonameid",
			"name",
			"name_i18n",
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
			$3,
			ST_GeomFromText($4),
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15
		)
	`,
		tmpRow.geonameid,
		fixName(tmpRow.asciiname),
		names,
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
