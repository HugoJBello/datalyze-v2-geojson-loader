package loaders

import (
	"database/sql"
	"datalyze-v2-geojson-loader-postgis/db"
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq" // here
)

func LoadPropertiesGeojson(jsonFile *os.File) (err error) {
	db.Connect()
	dbConnection := db.GetDb()

	if jsonFile == nil {
		return errors.New("no file")
	}
	geometry := models.GeojsonFromFile(jsonFile)
	properties := geometry.PropertyNames()

	tableName := "example"

	err = createDBwithProperties(dbConnection, properties, tableName)
	if err != nil {
		return err
		fmt.Println(err)
	}

	err = insertGeojsonWithProperties(dbConnection, geometry, tableName)
	if err != nil {
		return err
		fmt.Println(err)
	}
	return nil

}

func createDBwithProperties(dbConnection *sql.DB, properties []interface{}, tableName string) error {
	fmt.Println("creating tables")
	sql := fmt.Sprintf(`set client_encoding to 'utf8';
	DROP TABLE if exists %q;
	CREATE EXTENSION IF NOT EXISTS postgis;
	CREATE TABLE IF NOT EXISTS  %q (
		id SERIAL PRIMARY KEY,
		geom geometry(Multipolygon, 3857) NOT NULL
	);`, tableName, tableName)
	fmt.Println(sql)
	_, err := dbConnection.Exec(sql)

	return err
}

// readGeoJSON demonstrates reading data in GeoJSON format and inserting it
// into a database in EWKB format.
func insertGeojsonWithProperties(dbConnection *sql.DB, geometry models.Geojson, tableName string) error {

	for _, feature := range geometry.Features {
		err := insertFeatureWithProperties(dbConnection, feature, tableName)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil

}

func insertFeatureWithProperties(dbConnection *sql.DB, feature models.Feature, tableName string) error {

	b, err := json.Marshal(feature.Geometry)
	if err != nil {
		fmt.Println(err)
	}
	json := string(b)
	sql := fmt.Sprintf(`INSERT INTO %q (geom) VALUES (ST_SetSRID(ST_GeomFromGeoJSON($1), 3857));`, tableName)

	_, err = dbConnection.Exec(sql, json)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
