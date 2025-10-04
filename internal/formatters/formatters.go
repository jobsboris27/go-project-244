package formatters

import (
	models "code/internal/models"
)

const (
	STYLISH = "stylish"
	PLAIN   = "plain"
	JSON    = "json"
)

func RenderWithFormat(diffNodes []*models.DiffNode, format string) string {
	switch format {
	case PLAIN:
		return RenderPlain(diffNodes, "")
	case JSON:
		return RenderJSON(diffNodes)
	case STYLISH:
		return RenderStylish(diffNodes, 0)
	default:
		return RenderStylish(diffNodes, 0)
	}
}
