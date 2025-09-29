package formatters

type DiffNode struct {
	Key      string
	Status   string
	OldValue interface{}
	NewValue interface{}
	Children []*DiffNode
}

func RenderWithFormat(diffNodes []*DiffNode, format string) string {
	switch format {
	case "plain":
		return RenderPlain(diffNodes, "")
	case "json":
		return RenderJSON(diffNodes)
	case "stylish":
		return RenderStylish(diffNodes, 0)
	default:
		return RenderStylish(diffNodes, 0)
	}
}