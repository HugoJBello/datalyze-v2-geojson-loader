package loaders

import (
	"database/sql"
	"datalyze-v2-geojson-loader-postgis/db"
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/lib/pq" // here
)

func LoadPropertiesGeojson(jsonFile *os.File) (err error) {
	tableName := getTableName(jsonFile)

	db.Connect()
	dbConnection := db.GetDb()

	if jsonFile == nil {
		return errors.New("no file")
	}
	geometry := models.GeojsonFromFile(jsonFile)
	properties := geometry.PropertyNames()

	err = createDBwithProperties(dbConnection, properties, tableName)
	if err != nil {
		return err
		fmt.Println(err)
	}

	err = insertGeojsonWithProperties(dbConnection, geometry, tableName, properties)
	if err != nil {
		return err
		fmt.Println(err)
	}
	return nil

}

func getTableName(jsonFile *os.File) string {
	dirName := filepath.Dir(jsonFile.Name())
	tableName := strings.ReplaceAll(jsonFile.Name(), dirName, "")
	tableName = strings.ReplaceAll(tableName, ".json", "")
	tableName = strings.ReplaceAll(tableName, "/", "")
	return tableName
}

func createDBwithProperties(dbConnection *sql.DB, properties map[string]string, tableName string) error {
	fmt.Println("creating tables")
	propertiesString := generateColumnsFromProperties(properties)

	sql := fmt.Sprintf(`set client_encoding to 'utf8';
	DROP TABLE if exists %q;
	CREATE EXTENSION IF NOT EXISTS postgis;
	CREATE TABLE IF NOT EXISTS  %q (
		id SERIAL PRIMARY KEY,
		%q
		geom geometry(Multipolygon, 3857) NOT NULL
	);`, tableName, tableName, propertiesString)
	sql = strings.ReplaceAll(sql, "\"", "")
	sql = strings.ReplaceAll(sql, "%", "percent")

	_, err := dbConnection.Exec(sql)

	return err
}

func generateColumnsFromProperties(properties map[string]string) string {
	columns := ""

	for property := range properties {
		columns = columns + " " + property + " " + translateToSqlType(properties[property]) + ","
	}
	return columns
}

func translateToSqlType(golangType string) (sqlType string) {
	if golangType == "string" {
		return "TEXT"
	} else if strings.Contains(golangType, "float") {
		return "FLOAT"
	} else {
		return "INT"
	}
}

// readGeoJSON demonstrates reading data in GeoJSON format and inserting it
// into a database in EWKB format.
func insertGeojsonWithProperties(dbConnection *sql.DB, geometry models.Geojson, tableName string, properties map[string]string) error {

	for _, feature := range geometry.Features {
		err := insertFeatureWithProperties(dbConnection, feature, tableName, properties)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil

}

// https://stackoverflow.com/questions/17555857/go-unpacking-array-as-arguments
func insertFeatureWithProperties(dbConnection *sql.DB, feature models.Feature, tableName string, properties map[string]string) error {

	b, err := json.Marshal(feature.Geometry)
	if err != nil {
		fmt.Println(err)
	}
	json := string(b)
	propertiesText, propertiesSlice := generateColumnsInsertionFromProperties(properties)
	sql := fmt.Sprintf(`INSERT INTO %q (%q, geom) VALUES (%q ST_SetSRID(ST_GeomFromGeoJSON($1), 3857));`, tableName, propertiesText, generateValuesInsertionFromProperties(properties))
	sql = strings.ReplaceAll(sql, "\"", "")
	sql = strings.ReplaceAll(sql, "%", "percent")
	var args []reflect.Value
	args = generateArgs(sql, json, feature, propertiesSlice)
	fun := reflect.ValueOf(dbConnection.Exec)
	result := fun.Call(args)

	if result != nil && result[1].Interface() != nil {
		fmt.Println(result[1].Interface().(error))
	}
	return nil
}

func generateColumnsInsertionFromProperties(properties map[string]string) (string, []string) {
	columns := ""
	var columnsSlice []string

	for property := range properties {
		columnsSlice = append(columnsSlice, property)
		if columns != "" {
			columns = columns + "," + property
		} else {
			columns = columns + " " + property
		}
	}
	return columns, columnsSlice
}

func generateValuesInsertionFromProperties(properties map[string]string) string {
	values := ""
	count := 2
	for _ = range properties {
		values = values + fmt.Sprintf("$%d,", count)
		count = count + 1
	}
	return values
}

func generateArgs(sql string, json string, feature models.Feature, propertiesSlice []string) (args []reflect.Value) {
	args = append(args, reflect.ValueOf(sql))
	args = append(args, reflect.ValueOf(json))

	for _, value := range propertiesSlice {
		if feature.Properties[value] == nil {
			args = append(args, reflect.ValueOf(""))
		} else {
			args = append(args, reflect.ValueOf(feature.Properties[value]))
		}
	}
	return args
}
