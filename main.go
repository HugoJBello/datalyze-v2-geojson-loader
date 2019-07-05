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
//go run . -generate-from-csv -path data/csv_data/example.csv -output data/processed_data
// go run . -generate-from-csv -path data/csv_data/example.csv -output data/processed_data -output-mapbox
// go run . -generate-index

var jsonPath = flag.String("path", "data/example.json", "json path file")
var output = flag.String("output", "data/processed_data", "output dir")

var generateIndex = flag.Bool("generate-index", false, "load index into file")
var asMapbox = flag.Bool("output-mapbox", false, "create a mapbox tile (requires docker)")
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
		utils.GenerateMunicipioIndex()
		utils.GenerateCusecIndex()
	} else if *generateFromCsv {
		generators.GenerateGeojsonFromCsv(inputFile, outputPath)
		if *asMapbox {
			convetToMapbox(path, outputPath)
		}
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

func convetToMapbox(inputFile string, outputPath string) error {
	filenameInput := strings.ReplaceAll(inputFile, filepath.Dir(inputFile), "")
	outputPath = strings.ReplaceAll(outputPath+filenameInput, ".csv", "_geojson.json")
	outputMapbox := strings.ReplaceAll(outputPath, "_geojson.json", ".mbtiles")
	fmt.Println("converting " + outputPath + " into " + outputMapbox)
	fmt.Println("run the following command:")
	script := "docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o " + outputMapbox + " " + outputPath
	fmt.Println(script)
	out, err := exec.Command("docker", "run", "-it", "-v", "$(pwd):/data", "fuzzytolerance/tippecanoe", "tippecanoe", "-z12", "-f", "-o", outputMapbox, outputPath).Output()
	if err != nil {
		fmt.Println("docker error:", err)
	} else {
		text := string(out)
		fmt.Println(text)
	}
	fmt.Printf("%s", out)

	return nil
}
