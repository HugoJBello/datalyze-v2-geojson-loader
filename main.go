package main

import (
	"datalyze-v2-geojson-loader-postgis/generators"
	"datalyze-v2-geojson-loader-postgis/loaders"
	"datalyze-v2-geojson-loader-postgis/utils"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go run . -load-geojson -path data/raw_data/polygon_cusecs_data.tar.xz

var jsonPath = flag.String("path", "data/example.json", "json path file")
var output = flag.String("o", "data/processed_data", "output dir")

var generateIndex = flag.Bool("generate-index", false, "load index into file")
var generateFromCsv = flag.Bool("generate-from-csv", false, "generate geojson from csv")
var loadGeojson = flag.Bool("load-geojson", false, "load geojson in postgres")

func main() {
	flag.Parse()

	path := *jsonPath
	var inputFile *os.File

	if path != "" {
		inputFile = uncompressAndOpen(path)
	}

	outputPath := *output

	if *loadGeojson {
		loaders.LoadRawGeojson(inputFile)
	} else if *generateIndex {
		utils.GenerateCusecIndex()
	} else if *generateFromCsv {
		generators.GenerateGeojsonFromCsv(inputFile, outputPath)
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
		uncompressedinputFile, _ := os.Open(jsonPath)
		return uncompressedinputFile
	}
	inputFile, _ := os.Open(path)
	return inputFile
}
