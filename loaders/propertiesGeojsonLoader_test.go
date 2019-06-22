package loaders

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadWithProperties(t *testing.T) {
	testFile, _ := os.Open("../data/raw_data/polygon_cusecs_data.json")

	err := LoadPropertiesGeojson(testFile)
	assert.Equal(t, nil, err, "OK response is expected")

}
