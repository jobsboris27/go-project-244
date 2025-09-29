package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type TreeNode struct {
	Key      string
	Value    interface{}
	Children []*TreeNode
}

type DiffNode struct {
	Key      string
	Status   string
	OldValue interface{}
	NewValue interface{}
	Children []*DiffNode
}

const (
	YAML_EXT = ".yaml"
	JSON_EXT = ".json"
)

func Parse(path1, path2 string) string {
	return ParseWithFormat(path1, path2, "stylish")
}

func ParseWithFormat(path1, path2, format string) string {
	data1, err := parseByExtension(path1)
	if err != nil {
		fmt.Println("Error parsing file 1:", err)
		return ""
	}

	data2, err := parseByExtension(path2)
	if err != nil {
		fmt.Println("Error parsing file 2:", err)
		return ""
	}

	diff := genDiff(convertMapToTree(data1), convertMapToTree(data2))
	return renderWithFormat(diff, format)
}

func renderWithFormat(diffNodes []*DiffNode, format string) string {
	switch format {
	case "stylish":
		return renderStylish(diffNodes, 0)
	default:
		return renderStylish(diffNodes, 0)
	}
}

func parseByExtension(path string) (map[string]interface{}, error) {
	ext := filepath.Ext(path)
	switch ext {
	case JSON_EXT:
		return parseJSON(path)
	case YAML_EXT:
		return parseYAML(path)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}

func parseJSON(path string) (map[string]interface{}, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	errJson := json.Unmarshal(file, &result)

	if errJson != nil {
		return nil, errJson
	}
	return result, nil
}

func parseYAML(path string) (map[string]interface{}, error) {
	var result map[string]interface{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}

func convertMapToTree(data map[string]interface{}) *TreeNode {
	root := &TreeNode{Key: "root", Children: []*TreeNode{}}

	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			childTree := convertMapToTree(v)
			childNode := &TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case map[interface{}]interface{}:
			// Convert map[interface{}]interface{} to map[string]interface{}
			convertedMap := make(map[string]interface{})
			for k, val := range v {
				if strKey, ok := k.(string); ok {
					convertedMap[strKey] = val
				}
			}
			childTree := convertMapToTree(convertedMap)
			childNode := &TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case []interface{}:
			arrayNode := &TreeNode{
				Key:      key,
				Value:    nil,
				Children: []*TreeNode{},
			}

			for i, item := range v {
				itemNode := &TreeNode{
					Key:   fmt.Sprintf("[%d]", i),
					Value: item,
				}

				switch v := item.(type) {
				case map[string]interface{}:
					itemNode = convertMapToTree(v)
					itemNode.Key = fmt.Sprintf("[%d]", i)
				case []interface{}:
					itemNode = convertMapToTree(map[string]interface{}{
						fmt.Sprintf("[%d]", i): v,
					})
				default:
					itemNode = &TreeNode{
						Key:   fmt.Sprintf("[%d]", i),
						Value: item,
					}
				}

				arrayNode.Children = append(arrayNode.Children, itemNode)
			}
			root.Children = append(root.Children, arrayNode)

		default:
			leafNode := &TreeNode{
				Key:   key,
				Value: v,
			}
			root.Children = append(root.Children, leafNode)
		}
	}

	return root
}

func genDiff(tree1, tree2 *TreeNode) []*DiffNode {
	var diff []*DiffNode

	allKeys := collectAllKeys(tree1, tree2)
	sort.Strings(allKeys)

	for _, key := range allKeys {
		node1 := findChildByKey(tree1, key)
		node2 := findChildByKey(tree2, key)

		diffNode := &DiffNode{Key: key}

		switch {
		case node1 == nil && node2 != nil:
			diffNode.Status = "added"
			if hasChildren(node2) {
				diffNode.NewValue = reconstructObject(node2)
			} else {
				diffNode.NewValue = node2.Value
			}

		case node1 != nil && node2 == nil:
			diffNode.Status = "removed"
			if hasChildren(node1) {
				diffNode.OldValue = reconstructObject(node1)
			} else {
				diffNode.OldValue = node1.Value
			}

		case node1 != nil && node2 != nil:
			if hasChildren(node1) && hasChildren(node2) {
				diffNode.Status = "nested"
				diffNode.Children = genDiff(node1, node2)
			} else if !hasChildren(node1) && !hasChildren(node2) && areValuesEqual(node1, node2) {
				diffNode.Status = "unchanged"
				diffNode.OldValue = node1.Value
			} else {
				diffNode.Status = "modified"
				if hasChildren(node1) {
					diffNode.OldValue = reconstructObject(node1)
				} else {
					diffNode.OldValue = node1.Value
				}
				if hasChildren(node2) {
					diffNode.NewValue = reconstructObject(node2)
				} else {
					diffNode.NewValue = node2.Value
				}
			}
		}

		diff = append(diff, diffNode)
	}

	return diff
}

func renderStylish(diffNodes []*DiffNode, depth int) string {
	var result strings.Builder

	if depth == 0 {
		result.WriteString("{\n")
	}

	for _, node := range diffNodes {
		switch node.Status {
		case "unchanged":
			indent := strings.Repeat(" ", depth*4+4)
			result.WriteString(fmt.Sprintf("%s%s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))

		case "added":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s+ %s: %s\n", indent, node.Key, formatValue(node.NewValue, depth+1)))

		case "removed":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))

		case "modified":
			indent := strings.Repeat(" ", depth*4+2)
			result.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, node.Key, formatValue(node.OldValue, depth+1)))
			result.WriteString(fmt.Sprintf("%s+ %s: %s\n", indent, node.Key, formatValue(node.NewValue, depth+1)))

		case "nested":
			indent := strings.Repeat(" ", depth*4+4)
			result.WriteString(fmt.Sprintf("%s%s: {\n", indent, node.Key))
			result.WriteString(renderStylish(node.Children, depth+1))
			result.WriteString(fmt.Sprintf("%s}\n", indent))
		}
	}

	if depth == 0 {
		result.WriteString("}")
	}

	return result.String()
}

func formatValue(value interface{}, depth int) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return ""
		}
		return v
	case map[string]interface{}:
		return formatObject(v, depth)
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, val := range v {
			if strKey, ok := k.(string); ok {
				converted[strKey] = val
			}
		}
		return formatObject(converted, depth)
	case []interface{}:
		return formatArray(v, depth)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatObject(obj map[string]interface{}, depth int) string {
	if len(obj) == 0 {
		return "{}"
	}

	var result strings.Builder
	result.WriteString("{\n")

	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		indent := strings.Repeat(" ", depth*4+4)
		result.WriteString(fmt.Sprintf("%s%s: %s\n", indent, key, formatValue(obj[key], depth+1)))
	}

	indent := strings.Repeat(" ", depth*4)
	result.WriteString(fmt.Sprintf("%s}", indent))
	return result.String()
}

func formatArray(arr []interface{}, depth int) string {
	if len(arr) == 0 {
		return "[]"
	}

	var result strings.Builder
	result.WriteString("[\n")

	for i, item := range arr {
		indent := strings.Repeat(" ", depth*4+4)
		result.WriteString(fmt.Sprintf("%s%s\n", indent, formatValue(item, depth+1)))
		if i < len(arr)-1 {
			result.WriteString(",")
		}
	}

	indent := strings.Repeat(" ", depth*4)
	result.WriteString(fmt.Sprintf("%s]", indent))
	return result.String()
}

func collectAllKeys(tree1, tree2 *TreeNode) []string {
	keys := make(map[string]bool)

	if tree1 != nil {
		for _, child := range tree1.Children {
			keys[child.Key] = true
		}
	}
	if tree2 != nil {
		for _, child := range tree2.Children {
			keys[child.Key] = true
		}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}

func findChildByKey(tree *TreeNode, key string) *TreeNode {
	if tree == nil {
		return nil
	}
	for _, child := range tree.Children {
		if child.Key == key {
			return child
		}
	}
	return nil
}

func areValuesEqual(node1, node2 *TreeNode) bool {
	if node1 == nil || node2 == nil {
		return false
	}
	return fmt.Sprintf("%v", node1.Value) == fmt.Sprintf("%v", node2.Value)
}

func hasChildren(node *TreeNode) bool {
	return node != nil && len(node.Children) > 0
}

func reconstructObject(node *TreeNode) interface{} {
	if !hasChildren(node) {
		return node.Value
	}

	result := make(map[string]interface{})
	for _, child := range node.Children {
		result[child.Key] = reconstructObject(child)
	}
	return result
}
