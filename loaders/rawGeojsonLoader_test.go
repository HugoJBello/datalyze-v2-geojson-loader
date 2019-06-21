package loaders

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRaw(t *testing.T) {
	testFile, _ := os.Open("../data/raw_data/example.json")

	err := LoadRawGeojson(testFile)
	assert.Equal(t, nil, err, "OK response is expected")

}
