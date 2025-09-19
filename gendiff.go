package code

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type TreeNode struct {
	Key      string
	Value    interface{}
	Children []*TreeNode
}

const (
	YAML_EXT = ".yaml"
	JSON_EXT = ".json"
)

func Parse(path1, path2 string) {
	data1, err := parseByExtension(path1)
	if err != nil {
		fmt.Println("Error parsing file 1:", err)
		return
	}

	data2, err := parseByExtension(path2)
	if err != nil {
		fmt.Println("Error parsing file 2:", err)
		return
	}

	fmt.Println("File 1 data:", convertMapToTree(data1))
	fmt.Println("File 2 data:", convertMapToTree(data2))
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
			childNode := convertMapToTree(v)
			childNode.Key = key
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
