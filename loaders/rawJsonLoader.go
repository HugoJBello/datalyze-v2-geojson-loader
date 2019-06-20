package loaders

import (
	"database/sql"
	"datalyze-v2-geojson-loader-postgis/db"
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq" // here
)

func LoadRawGeojson(jsonFile *os.File) {
	db.Connect()
	dbConnection := db.GetDb()
	err := createDBRawGeojson(dbConnection)
	if err != nil {
		fmt.Println(err)
	}

	geometry := models.GeojsonFromFile(jsonFile)
	// fmt.Println(geometry)
	err = insertRawGeojson(dbConnection, geometry)
	if err != nil {
		fmt.Println(err)
	}

}

func createDBRawGeojson(dbConnection *sql.DB) error {
	fmt.Println("creating tables")
	_, err := dbConnection.Exec(`
	set client_encoding to 'utf8';
		DROP TABLE waypoints;
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS waypoints (
			id SERIAL PRIMARY KEY,
			geom geometry(Multipolygon, 3857) NOT NULL
		);
	`)
	return err
}

// readGeoJSON demonstrates reading data in GeoJSON format and inserting it
// into a database in EWKB format.
func insertRawGeojson(dbConnection *sql.DB, geometry models.Geojson) error {

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

	b, err := json.Marshal(feature.Geometry)
	if err != nil {
		fmt.Println(err)
	}
	json := string(b)

	_, err = dbConnection.Exec(`
	INSERT INTO waypoints(geom) VALUES (ST_SetSRID(ST_GeomFromGeoJSON($1), 3857));
	`, json)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
