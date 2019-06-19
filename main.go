package main

import (
	"database/sql"
	"datalyze-v2-geojson-loader-postgis/db"
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq" // here
)

func main() {
	db.Connect()
	dbConnection := db.GetDb()
	err := createDB(dbConnection)
	if err != nil {
		fmt.Println(err)
	}

	geometry := geojsonFromFile("example.json")
	err = insertGeojson(dbConnection, geometry)
	if err != nil {
		fmt.Println(err)
	}

}

func createDB(dbConnection *sql.DB) error {
	fmt.Println("creating tables")
	_, err := dbConnection.Exec(`
	set client_encoding to 'utf8';
		DROP TABLE waypoints;
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS waypoints (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			cusec TEXT NOT NULL,
			percent_pob_esp FLOAT,
			geom geometry(Multipolygon) NOT NULL
		);
	`)
	return err
}

func geojsonFromFile(filename string) (geometry models.Geojson) {

	jsonFile, _ := os.Open(filename)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	geometry = models.Geojson{}
	json.Unmarshal(byteValue, &geometry)

	return geometry
}

// readGeoJSON demonstrates reading data in GeoJSON format and inserting it
// into a database in EWKB format.
func insertGeojson(dbConnection *sql.DB, geometry models.Geojson) error {

	for _, feature := range geometry.Features {
		err := insertFeature(dbConnection, feature)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil

}

func insertFeature(dbConnection *sql.DB, feature models.Feature) error {
	percentPobEsp := feature.Properties["%_pob_esp"]
	cusec := feature.Properties["CUSEC"]

	b, err := json.Marshal(feature.Geometry)
	if err != nil {
		fmt.Println(err)
	}
	json := string(b)

	_, err = dbConnection.Exec(`
	INSERT INTO waypoints(name, cusec, percent_pob_esp, geom) VALUES ($1, $2, $3, ST_GeomFromGeoJSON($4));
	`, "aa", cusec, percentPobEsp, json)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
