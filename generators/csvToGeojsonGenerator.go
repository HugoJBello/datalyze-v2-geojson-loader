package generators

import (
	"bufio"
	"datalyze-v2-geojson-loader-postgis/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func GenerateGeojsonFromCsv(file *os.File) (geojsonResult models.Geojson, err error) {
	cusecIndex := getCusecIndex()
	geojsonResult = models.Geojson{}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	var header []string
	count := 1
	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		if count == 1 {
			header = line
			count = count + 1
		} else {
			properties := obtainJsonFromLine(line, header)
			feature := cusecIndex[fmt.Sprintf("%v", properties["CUSEC"])]
			feature.Properties = properties
			geojsonResult.Features = append(geojsonResult.Features, feature)
		}
	}
	return geojsonResult, nil
}

func obtainJsonFromLine(line []string, header []string) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	for index, column := range header {
		jsonMap[column] = line[index]
	}

	return jsonMap
}

func getCusecIndex() map[string]models.Feature {
	cusecIndexFile, err := os.Open("../data/cusec_index.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(cusecIndexFile)

	cusecIndex := make(map[string]models.Feature)
	json.Unmarshal(byteValue, &cusecIndex)
	return cusecIndex
}
