package utils

import (
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/json"
	"fmt"
	"os"
)

func GenerateMunicipioIndex() {
	jsonFile := uncompressAndOpen("./data/raw_data/municipios_epsg4326_CRS84.geojson.tar.xz")
	//jsonFile := uncompressAndOpen("../data/raw_data/example2_geojson.tar.xz")

	geojson := models.GeojsonFromFile(jsonFile)

	mapCusecs := make(map[string]models.Feature)

	fmt.Println("creating municipio index")
	for _, feature := range geojson.Features {
		if feature.Properties["CODIGOINE"] != nil {
			mapCusecs[fmt.Sprintf("%v", feature.Properties["CODIGOINE"])] = feature
		}
	}
	jsonString, err := json.Marshal(mapCusecs)
	if err != nil {
		fmt.Println(err)
	}

	f, _ := os.Create("./data/municipio_index.json")
	_, err = f.WriteString(string(jsonString))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

}
