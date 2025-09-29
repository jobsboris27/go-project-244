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

func TestParseFlatYAML(t *testing.T) {
	result := Parse("testdata/fixture/file1.yaml", "testdata/fixture/file2.yaml")

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`

	expectedStr := strings.TrimSpace(expected)
	resultStr := strings.TrimSpace(result)

	assert.Equal(t, expectedStr, resultStr)
}

func TestParseYAMLFile(t *testing.T) {
	data, err := parseYAML("testdata/fixture/file1.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	assert.Equal(t, "hexlet.io", data["host"])
	assert.Equal(t, 50, data["timeout"])
	assert.Equal(t, false, data["follow"])
}


func TestParseNestedJSON(t *testing.T) {
	result := Parse("testdata/fixture/nested1.json", "testdata/fixture/nested2.json")

	assert.Contains(t, result, "common: {")
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "- setting2: 200")
	assert.Contains(t, result, "- setting3: true")
	assert.Contains(t, result, "+ setting3: null")
	assert.Contains(t, result, "+ setting5: {")
	assert.Contains(t, result, "key5: value5")
	assert.Contains(t, result, "- group2: {")
	assert.Contains(t, result, "+ group3: {")
}

func TestParseNestedYAML(t *testing.T) {
	result := Parse("testdata/fixture/nested1.yaml", "testdata/fixture/nested2.yaml")

	assert.Contains(t, result, "common: {")
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "- setting2: 200")
	assert.Contains(t, result, "- setting3: true")
	assert.Contains(t, result, "+ setting3: null")
}

func TestParseWithFormat(t *testing.T) {
	result1 := ParseWithFormat("testdata/fixture/nested1.json", "testdata/fixture/nested2.json", "stylish")
	result2 := Parse("testdata/fixture/nested1.json", "testdata/fixture/nested2.json")

	assert.Equal(t, result1, result2)
}

