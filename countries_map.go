package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type tmpCountryInfoIDToISO struct {
	geonameid string
	iso       string
}

type tmpShapesAllLow struct {
	polygonData string
}

func listOfNeighbours(lst string) []int {
	lsta := strings.Split(lst, ",")
	result := []int{}
	for _, isoName := range lsta {
		if _, ok := isoToGeonameID[isoName]; ok {
			result = append(result, isoToGeonameID[isoName])
		}
	}
	return result
}

func getCountryPolygon(db *sql.DB, geoname int) string {
	rows, err := db.Query(`
		SELECT
			"polygonData"
		FROM "tmp_shapesAllLow"
		WHERE "geonameid_1" = $1
	`, geoname)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r tmpShapesAllLow
		e := rows.Scan(
			&r.polygonData,
		)
		if e == nil {
			return r.polygonData
		} else {
			log.Fatal(err)
		}
	}
	return ""
}

func getIsoToGeonameID(db *sql.DB) (map[string]int, map[int]string) {
	isoToGeonameIDL := make(map[string]int)
	geonameIDToISOL := make(map[int]string)
	rows, err := db.Query(`
		SELECT
			"geonameid",
			"iso"
		FROM "tmp_countryInfo"
		WHERE "geonameid" != ''
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r tmpCountryInfoIDToISO
		e := rows.Scan(
			&r.geonameid,
			&r.iso,
		)
		if e == nil {
			id := toInt(r.geonameid)
			isoToGeonameIDL[r.iso] = id
			geonameIDToISOL[id] = r.iso
		} else {
			log.Fatal(err)
		}
	}
	return isoToGeonameIDL, geonameIDToISOL
}
