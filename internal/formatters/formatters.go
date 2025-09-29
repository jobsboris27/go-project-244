package formatters

import (
	models "code/internal/models"
)

func RenderWithFormat(diffNodes []*models.DiffNode, format string) string {
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
