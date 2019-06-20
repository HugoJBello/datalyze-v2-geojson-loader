package main

import (
	"datalyze-v2-geojson-loader-postgis/loaders"
	"flag"
	"os"
)

var jsonPath = flag.String("path", "data/example.json", "json path file")
var loadGeojson = flag.Bool("load-geojson", false, "load geojson in postgres")

func main() {
	flag.Parse()

	path := *jsonPath
	jsonFile, _ := os.Open(path)
	if *loadGeojson {
		loaders.LoadRawGeojson(jsonFile)
	}
}
