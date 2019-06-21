package main

import (
	"datalyze-v2-geojson-loader-postgis/loaders"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var jsonPath = flag.String("path", "data/example.json", "json path file")
var loadGeojson = flag.Bool("load-geojson", false, "load geojson in postgres")

func main() {
	flag.Parse()

	path := *jsonPath

	jsonFile := uncompressAndOpen(path)

	if *loadGeojson {
		loaders.LoadRawGeojson(jsonFile)
	}
}

func uncompressAndOpen(path string) *os.File {
	if strings.HasSuffix(path, "tar.xz") {
		fmt.Println("Uncompresing geojson inside " + path)
		outputDir := filepath.Dir(path)

		cmd := exec.Command("tar", "-x", "-f", path, "-C", outputDir)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		jsonPath := strings.ReplaceAll(path, ".tar.xz", ".json")
		uncompressedJsonFile, _ := os.Open(jsonPath)
		return uncompressedJsonFile
	}
	jsonFile, _ := os.Open(path)
	return jsonFile
}
