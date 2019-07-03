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
	"path/filepath"
	"strconv"
	"strings"
)

var numberNameFileds []string = []string{"RANKING", "MEDIANA", "PERCENT", "percent", "PSOE", "kpi", "KPI"}

func GenerateGeojsonFromCsv(file *os.File, outputPath string) (err error) {
	fmt.Println("Reading previously obtained cusec index")
	cusecIndex := getCusecIndex()
	filename := strings.ReplaceAll(file.Name(), filepath.Dir(file.Name()), "")

	geojsonResult := initializeGeojsonWithCs(filename)

	fmt.Println("Reading input file")
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
			return err
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
	filenameOut := strings.ReplaceAll(outputPath+filename, ".csv", "_geojson.json")
	filenameOut = strings.ReplaceAll(filenameOut, "//", "/")

	fmt.Println("saving file " + filenameOut)
	saveToJsonFile(geojsonResult, filenameOut)
	return nil
}

func initializeGeojsonWithCs(filename string) models.Geojson {
	geojsonResult := models.Geojson{}
	geojsonResult.Name = strings.ReplaceAll(filename, "/", "")
	geojsonResult.Type = "FeatureCollection"
	geojsonResult.CRS = make(map[string]interface{})
	geojsonResult.CRS["type"] = "name"
	properties := make(map[string]interface{})
	properties["name"] = "urn:ogc:def:crs:OGC:1.3:CRS84"
	geojsonResult.CRS["properties"] = properties
	return geojsonResult
}

func obtainJsonFromLine(line []string, header []string) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	for index, column := range header {
		if !isColumnNumber(column) {
			jsonMap[column] = line[index]
		} else {
			jsonMap[column], _ = strconv.ParseFloat(line[index], 32)
		}
	}

	return jsonMap
}

func isColumnNumber(column string) bool {
	for _, name := range numberNameFileds {
		if strings.Contains(column, name) {
			return true
		}
	}
	return false
}

func getCusecIndex() map[string]models.Feature {
	cusecIndexFile, err := os.Open("../data/cusec_index.json")
	if err != nil {
		cusecIndexFile, _ = os.Open("./data/cusec_index.json")
	}
	byteValue, _ := ioutil.ReadAll(cusecIndexFile)

	cusecIndex := make(map[string]models.Feature)
	json.Unmarshal(byteValue, &cusecIndex)
	return cusecIndex
}

func saveToJsonFile(data interface{}, filename string) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(filename)

	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	_, err = f.WriteString(string(jsonString))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

}
