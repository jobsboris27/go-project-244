package code

import (
	"code/internal/formatters"
	parsers "code/internal/parsers"
	"fmt"
)

func GenDiff(path1, path2, format string) (string, error) {
	data1, err := parsers.ParseByExtension(path1)
	if err != nil {
		return "", fmt.Errorf("parsing file %s: %w", path1, err)
	}

	data2, err := parsers.ParseByExtension(path2)
	if err != nil {
		return "", fmt.Errorf("parsing file %s: %w", path1, err)
	}

	diff := parsers.GetDiff(parsers.СonvertMapToTree(data1), parsers.СonvertMapToTree(data2))
	return formatters.RenderWithFormat(diff, format), nil
}
