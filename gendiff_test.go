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
* timeout: 20
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

func TestDebugDiff(t *testing.T) {
	data1, _ := parseByExtension("file1.yaml")
	data2, _ := parseByExtension("file2.yaml")

	tree1 := convertMapToTree(data1)
	tree2 := convertMapToTree(data2)

	diff := GenDiff(tree1, tree2)
	t.Logf("Diff nodes: %d", len(diff))
	for _, node := range diff {
		t.Logf("  Key: %s, Status: %s, OldValue: %+v, NewValue: %+v, Children: %d",
			node.Key, node.Status, node.OldValue, node.NewValue, len(node.Children))
	}
}

func TestDebugComplexYAML(t *testing.T) {
	// Test simple data first
	simpleData := map[string]interface{}{
		"host": "test.com",
		"config": map[string]interface{}{
			"timeout": 30,
			"enabled": true,
		},
	}

	t.Logf("SimpleData: %+v", simpleData)
	simpleTree := convertMapToTree(simpleData)
	t.Logf("SimpleTree children: %d", len(simpleTree.Children))
	for _, child := range simpleTree.Children {
		t.Logf("  Key: %s, HasChildren: %t, Children: %d, Value: %+v", child.Key, len(child.Children) > 0, len(child.Children), child.Value)
		for _, subchild := range child.Children {
			t.Logf("    SubKey: %s, HasChildren: %t, Value: %+v", subchild.Key, len(subchild.Children) > 0, subchild.Value)
		}
	}

	// Now test complex files
	data1, err := parseByExtension("file1.yaml")
	assert.NoError(t, err)
	data2, err := parseByExtension("file2.yaml")
	assert.NoError(t, err)

	t.Logf("Data1: %+v", data1)
	t.Logf("Data2: %+v", data2)

	// Check types
	for k, v := range data1 {
		t.Logf("Data1[%s] type: %T, value: %+v", k, v, v)
	}

	tree1 := convertMapToTree(data1)
	tree2 := convertMapToTree(data2)

	t.Logf("Tree1 children: %d", len(tree1.Children))
	for _, child := range tree1.Children {
		t.Logf("  Key: %s, HasChildren: %t, Children: %d, Value: %+v", child.Key, len(child.Children) > 0, len(child.Children), child.Value)
		for _, subchild := range child.Children {
			t.Logf("    SubKey: %s, HasChildren: %t, Value: %+v", subchild.Key, len(subchild.Children) > 0, subchild.Value)
		}
	}

	t.Logf("Tree2 children: %d", len(tree2.Children))
	for _, child := range tree2.Children {
		t.Logf("  Key: %s, HasChildren: %t, Children: %d, Value: %+v", child.Key, len(child.Children) > 0, len(child.Children), child.Value)
		for _, subchild := range child.Children {
			t.Logf("    SubKey: %s, HasChildren: %t, Value: %+v", subchild.Key, len(subchild.Children) > 0, subchild.Value)
		}
	}
}
