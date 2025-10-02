package code

import (
	"code/internal/formatters"
	models "code/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v2"
)

const (
	YAML_EXT       = ".yaml"
	YAML_EXT_SHORT = ".yml"
	JSON_EXT       = ".json"
)

func GenDiff(path1, path2, format string) (string, error) {
	data1, err := parseByExtension(path1)
	if err != nil {
		fmt.Println("Error parsing file 1:", err)
		return "", err
	}

	data2, err := parseByExtension(path2)
	if err != nil {
		fmt.Println("Error parsing file 2:", err)
		return "", err
	}

	diff := genDiff(convertMapToTree(data1), convertMapToTree(data2))
	return formatters.RenderWithFormat(diff, format), nil
}

func parseByExtension(path string) (map[string]interface{}, error) {
	ext := filepath.Ext(path)
	switch ext {
	case JSON_EXT:
		return parseJSON(path)
	case YAML_EXT, YAML_EXT_SHORT:
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

	var raw interface{}
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, err
	}

	switch v := raw.(type) {
	case map[string]interface{}:
		return v, nil
	case []interface{}:
		return map[string]interface{}{"root": v}, nil
	default:
		return map[string]interface{}{"root": v}, nil
	}
}

func parseYAML(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var raw interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	switch v := raw.(type) {
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, val := range v {
			if strKey, ok := k.(string); ok {
				converted[strKey] = val
			} else {
				converted[fmt.Sprintf("%v", k)] = val
			}
		}
		return converted, nil
	case map[string]interface{}:
		return v, nil
	case []interface{}:
		return map[string]interface{}{"root": v}, nil
	default:
		return map[string]interface{}{"root": v}, nil
	}
}

func convertMapToTree(data map[string]interface{}) *models.TreeNode {
	root := &models.TreeNode{Key: "root", Children: []*models.TreeNode{}}

	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			childTree := convertMapToTree(v)
			childNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case map[interface{}]interface{}:
			convertedMap := make(map[string]interface{})
			for k, val := range v {
				if strKey, ok := k.(string); ok {
					convertedMap[strKey] = val
				}
			}
			childTree := convertMapToTree(convertedMap)
			childNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: childTree.Children,
			}
			root.Children = append(root.Children, childNode)

		case []interface{}:
			arrayNode := &models.TreeNode{
				Key:      key,
				Value:    nil,
				Children: []*models.TreeNode{},
			}

			for i, item := range v {
				itemNode := &models.TreeNode{
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
					itemNode = &models.TreeNode{
						Key:   fmt.Sprintf("[%d]", i),
						Value: item,
					}
				}

				arrayNode.Children = append(arrayNode.Children, itemNode)
			}
			root.Children = append(root.Children, arrayNode)
		default:
			leafNode := &models.TreeNode{
				Key:   key,
				Value: v,
			}
			root.Children = append(root.Children, leafNode)
		}
	}

	return root
}

func genDiff(tree1, tree2 *models.TreeNode) []*models.DiffNode {
	var diff []*models.DiffNode

	allKeys := collectAllKeys(tree1, tree2)
	sort.Strings(allKeys)

	for _, key := range allKeys {
		node1 := findChildByKey(tree1, key)
		node2 := findChildByKey(tree2, key)

		diffNode := &models.DiffNode{Key: key}

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

func collectAllKeys(tree1, tree2 *models.TreeNode) []string {
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

func findChildByKey(tree *models.TreeNode, key string) *models.TreeNode {
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

func areValuesEqual(node1, node2 *models.TreeNode) bool {
	if node1 == nil || node2 == nil {
		return false
	}
	return fmt.Sprintf("%v", node1.Value) == fmt.Sprintf("%v", node2.Value)
}

func hasChildren(node *models.TreeNode) bool {
	return node != nil && len(node.Children) > 0
}

func reconstructObject(node *models.TreeNode) interface{} {
	if !hasChildren(node) {
		return node.Value
	}

	result := make(map[string]interface{})
	for _, child := range node.Children {
		result[child.Key] = reconstructObject(child)
	}
	return result
}
