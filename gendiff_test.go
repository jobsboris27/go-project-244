package code

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlatJSON(t *testing.T) {
	result, _ := Parse("testdata/fixture/file1.json", "testdata/fixture/file2.json")
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
	result, _ := Parse("testdata/fixture/file1.yaml", "testdata/fixture/file2.yaml")

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
	result, _ := Parse("testdata/fixture/nested1.json", "testdata/fixture/nested2.json")

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
	result, _ := Parse("testdata/fixture/nested1.yaml", "testdata/fixture/nested2.yaml")

	assert.Contains(t, result, "common: {")
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "- setting2: 200")
	assert.Contains(t, result, "- setting3: true")
	assert.Contains(t, result, "+ setting3: null")
}

func TestGenDiff(t *testing.T) {
	result1, _ := GenDiff("testdata/fixture/nested1.json", "testdata/fixture/nested2.json", "stylish")
	result2, _ := Parse("testdata/fixture/nested1.json", "testdata/fixture/nested2.json")

	assert.Equal(t, result1, result2)
}

func TestPlainFormatter(t *testing.T) {
	result, _ := GenDiff("testdata/fixture/nested1.json", "testdata/fixture/nested2.json", "plain")

	assert.Contains(t, result, "Property 'common.follow' was added with value: false")
	assert.Contains(t, result, "Property 'common.setting2' was removed")
	assert.Contains(t, result, "Property 'common.setting3' was updated. From true to null")
	assert.Contains(t, result, "Property 'common.setting4' was added with value: 'blah blah'")
	assert.Contains(t, result, "Property 'common.setting5' was added with value: [complex value]")
	assert.Contains(t, result, "Property 'common.setting6.doge.wow' was updated. From '' to 'so much'")
	assert.Contains(t, result, "Property 'common.setting6.ops' was added with value: 'vops'")
	assert.Contains(t, result, "Property 'group1.baz' was updated. From 'bas' to 'bars'")
	assert.Contains(t, result, "Property 'group1.nest' was updated. From [complex value] to 'str'")
	assert.Contains(t, result, "Property 'group2' was removed")
	assert.Contains(t, result, "Property 'group3' was added with value: [complex value]")
}

func TestPlainFormatterSimple(t *testing.T) {
	result, _ := GenDiff("testdata/fixture/file1.yaml", "testdata/fixture/file2.yaml", "plain")

	assert.Contains(t, result, "Property 'follow' was removed")
	assert.Contains(t, result, "Property 'proxy' was removed")
	assert.Contains(t, result, "Property 'timeout' was updated. From 50 to 20")
	assert.Contains(t, result, "Property 'verbose' was added with value: true")
}

func TestJSONFormatter(t *testing.T) {
	result, _ := GenDiff("testdata/fixture/file1.yaml", "testdata/fixture/file2.yaml", "json")

	assert.Contains(t, result, `"key": "follow"`)
	assert.Contains(t, result, `"type": "removed"`)
	assert.Contains(t, result, `"key": "timeout"`)
	assert.Contains(t, result, `"type": "updated"`)
	assert.Contains(t, result, `"oldValue": 50`)
	assert.Contains(t, result, `"newValue": 20`)
	assert.Contains(t, result, `"key": "verbose"`)
	assert.Contains(t, result, `"type": "added"`)

	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(result), &jsonData)
	assert.NoError(t, err)

	diffArray, exists := jsonData["diff"].([]interface{})
	assert.True(t, exists)
	assert.Greater(t, len(diffArray), 0)
}

func TestJSONFormatterNested(t *testing.T) {
	result, _ := GenDiff("testdata/fixture/nested1.json", "testdata/fixture/nested2.json", "json")

	assert.Contains(t, result, `"key": "common"`)
	assert.Contains(t, result, `"type": "nested"`)
	assert.Contains(t, result, `"children"`)
	assert.Contains(t, result, `"key": "setting3"`)
	assert.Contains(t, result, `"oldValue": true`)
	assert.Contains(t, result, `"newValue": null`)

	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(result), &jsonData)
	assert.NoError(t, err)

	diffArray, exists := jsonData["diff"].([]interface{})
	assert.True(t, exists)
	assert.Greater(t, len(diffArray), 0)
}

func TestParseJSONArray(t *testing.T) {
	data, err := parseByExtension("testdata/fixture/file1arr.json")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	rootVal, ok := data["root"]
	assert.True(t, ok, "expected 'root' key for top-level array")

	arr, ok := rootVal.([]interface{})
	assert.True(t, ok, "expected root to be []interface{}")
	assert.Len(t, arr, 1)

	obj, ok := arr[0].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, "hexlet.io", obj["host"])

	assert.Equal(t, float64(50), obj["timeout"])
	assert.Equal(t, "123.234.53.22", obj["proxy"])
	assert.Equal(t, false, obj["follow"])
}

func TestParseJSONCompareArrays(t *testing.T) {
	data1, err := parseByExtension("testdata/fixture/file1arr.json")
	assert.NoError(t, err)
	data2, err := parseByExtension("testdata/fixture/file2arr.json")
	assert.NoError(t, err)

	obj1, ok := data1["root"].([]interface{})[0].(map[string]interface{})
	assert.True(t, ok)
	obj2, ok := data2["root"].([]interface{})[0].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, float64(50), obj1["timeout"])
	assert.Equal(t, float64(20), obj2["timeout"])
	assert.Equal(t, true, obj2["verbose"])
	assert.Equal(t, "hexlet.io", obj2["host"])

	assert.Equal(t, "123.234.53.22", obj1["proxy"])
	_, ok = obj2["proxy"]
	assert.False(t, ok)
}
