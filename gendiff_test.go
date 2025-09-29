package code

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlatJSON(t *testing.T) {
	result := Parse("testdata/fixture/file1.json", "testdata/fixture/file2.json")
	expected, err := os.ReadFile("testdata/fixture/expected_flat.txt")
	assert.NoError(t, err)

	expectedStr := strings.TrimSpace(string(expected))
	resultStr := strings.TrimSpace(result)

	assert.Equal(t, expectedStr, resultStr)
}

func TestParseByExtensionJSON(t *testing.T) {
	data, err := parseByExtension("testdata/fixture/file1.json")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	assert.Equal(t, "hexlet.io", data["host"])
	assert.Equal(t, float64(50), data["timeout"])
	assert.Equal(t, false, data["follow"])
}

func TestParseByExtensionUnsupported(t *testing.T) {
	_, err := parseByExtension("test.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file extension")
}

func TestParseJSONNonExistentFile(t *testing.T) {
	_, err := parseJSON("nonexistent.json")
	assert.Error(t, err)
}
