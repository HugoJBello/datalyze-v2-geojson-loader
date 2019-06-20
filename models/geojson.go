package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Geometry struct {
	Type        string    `json:"type"`
	BoundingBox []float64 `json:"bbox,omitempty"`
	Coordinates [][][][]float64
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects are not currently supported
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`             // Coordinate Reference System Objects are not currently supported
	Properties map[string]interface{} `json:"properties,omitempty"` // Coordinate Reference System Objects are not currently supported
}

type Geojson struct {
	Type     string                 `json:"type"`
	CRS      map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects are not currently supported
	Name     string                 `json:"name"`
	Features []Feature              `json:"features"`
}

func GeojsonFromFile(jsonFile *os.File) (geometry Geojson) {

	byteValue, _ := ioutil.ReadAll(jsonFile)

	geometry = Geojson{}
	json.Unmarshal(byteValue, &geometry)

	return geometry
}
