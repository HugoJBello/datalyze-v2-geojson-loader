package utils

import (
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GenerateCusecIndex() {
	jsonFile := uncompressAndOpen("./data/raw_data/polygon_cusecs_geojson.tar.xz")
	//jsonFile := uncompressAndOpen("../data/raw_data/example2_geojson.tar.xz")

	geojson := models.GeojsonFromFile(jsonFile)

	mapCusecs := make(map[string]models.Feature)

	fmt.Println("creating cusec index")
	for _, feature := range geojson.Features {
		if feature.Properties["CUSEC"] != nil {
			mapCusecs[fmt.Sprintf("%v", feature.Properties["CUSEC"])] = feature
		}
	}
	jsonString, err := json.Marshal(mapCusecs)
	if err != nil {
		fmt.Println(err)
	}

	f, _ := os.Create("./data/cusec_index.json")
	_, err = f.WriteString(string(jsonString))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
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
