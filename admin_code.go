package main

import (
	"database/sql"
	"log"
)

type tmpAdmin1Codes struct {
	code      string
	asciiname string
	geonameid int
}

func adminCode(db *sql.DB) map[string]int {
	var adminCodeMap = make(map[string]int)
	_, err := db.Exec(`TRUNCATE "adminCode" CASCADE;`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`
		SELECT
			"code",
			"asciiname",
			"geonameid"
		FROM "tmp_admin1Codes"
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r tmpAdmin1Codes
		err := rows.Scan(
			&r.code,
			&r.asciiname,
			&r.geonameid,
		)
		adminCodeMap[r.code] = r.geonameid
		if err != nil {
			log.Fatal(err)
		}
		insertAdminCode(db, r.geonameid, r.asciiname)
	}

	return adminCodeMap
}

func insertAdminCode(
	db *sql.DB,
	geonameid int,
	name string,
) {
	_, err := db.Exec(`
		INSERT INTO "adminCode" (
			"id",
			"name"
		) VALUES (
			$1,
			$2
		)
	`,
		geonameid,
		fixName(name),
	)

	if err != nil {
		log.Fatal(err)
	}
}
