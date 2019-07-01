package generators

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvToGeojsonGenerator(t *testing.T) {
	testFile, err := os.Open("../data/csv_data/example.csv")
	fmt.Println(testFile)

	if err != nil {
		fmt.Println("error ")
		fmt.Println(err)
	}

	err = GenerateGeojsonFromCsv(testFile)
	assert.Equal(t, nil, err, "OK response is expected")

}
