package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	_ "github.com/lib/pq"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var nameFixReplace = regexp.MustCompile(`[^a-zA-Z0-9- ]`)

func fixName(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return nameFixReplace.ReplaceAllString(result, "")
}

func toInt(s string) int {
	id, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return id
}

func findPrefereLanguage(locales string) string {
	lst := strings.Split(locales, ",")
	if len(lst) == 0 {
		return ""
	}
	if len(lst[0]) < 2 {
		return ""
	}
	return lst[0][0:2]
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt64(s string) sql.NullInt64 {
	if len(s) == 0 || s == "" {
		return sql.NullInt64{}
	}
	num, e := strconv.Atoi(s)
	if e != nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(num),
		Valid: true,
	}
}

func NewNullFloatInt(s string) sql.NullInt64 {
	if len(s) == 0 || s == "" {
		return sql.NullInt64{}
	}
	num, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(math.RoundToEven(num)),
		Valid: true,
	}
}

func parseLocales(locales string) []string {
	result := []string{}
	lst := strings.Split(locales, ",")
	if len(lst) == 0 {
		return result
	}
	for _, ll := range lst {
		_, err := language.Parse(ll)
		if err == nil {
			result = append(result, ll)
		}
	}
	return result
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

var isoToGeonameID map[string]int

func getDb() (*sql.DB, bool) {
	connectionString := os.Getenv("POSTGRES_URI")
	if connectionString == "" {
		connectionString = "postgres://geonames:geonames@127.0.0.1/geonames?sslmode=disable"
	}
	db, err1 := sql.Open("postgres", connectionString)

	if err1 == nil {
		rows, err2 := db.Query(`
			SELECT
				COUNT(*)
			FROM "tmp_ready"
		`)
		if err2 != nil {
			fmt.Println(err2)
			db.Close()
			return nil, false
		}
		defer rows.Close()

		var num int = 0
		for rows.Next() {
			err3 := rows.Scan(&num)

			if err3 != nil {
				fmt.Println(err3)
				db.Close()
				return nil, false
			}
		}

		if num > 0 {
			return db, true
		}
	} else {

		fmt.Println(err1)
	}

	return nil, false
}

func dockerLogs() {
	out1, err1 := exec.Command("/usr/local/bin/docker ps -a").Output()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Printf("docker ps\n%s\n", out1)
	out2, err2 := exec.Command("/usr/local/bin/docker logs -n 10").Output()
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Printf("docker logs\n%s\n", out2)
}

func main() {
	isoToGeonameID = make(map[string]int)

	var db *sql.DB
	var okDB bool = false
	for !okDB {
		fmt.Println("Try connect to database")
		db, okDB = getDb()

		dockerLogs()
		time.Sleep(time.Second * 15)
	}

	isoToGeonameID, _ = getIsoToGeonameID(db)

	fmt.Println("Process countries")
	processCountryInfo(db)
	fmt.Println("Process geo data")
	buildGeoData(db)

	fmt.Println("Drop temporary tables")
	db.Exec(`DROP TABLE "tmp_shapesAllLow";`)
	db.Exec(`DROP TABLE "tmp_hierarchy";`)
	db.Exec(`DROP TABLE "tmp_countryInfo";`)
	db.Exec(`DROP TABLE "tmp_geonameid";`)
	db.Exec(`DROP TABLE "tmp_alternateNamesV2";`)
	db.Exec(`DROP TABLE "tmp_ready";`)

	fmt.Println("Process geo data")
	db.Exec(`VACUUM(FULL, ANALYZE) "countryInfo";`)
	db.Exec(`VACUUM(FULL, ANALYZE) "geo";`)
}
