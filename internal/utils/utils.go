package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/xuri/excelize/v2"
)

func AbsolutePathTest() z.Test[*string] {
	return z.Test[*string]{
		Func: func(path *string, ctx z.Ctx) {
			if !filepath.IsAbs(*path) {
				ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("Path '%s' is not absolute", *path)))
			}
		},
	}
}

func ParseRange(rangeStr string) (int, int, int, int, error) {

	re := regexp.MustCompile(`(\$?[A-Z]+\$?\d+)(:\$?[A-Z]+\$?\d+)?`)
	matches := re.FindStringSubmatch(rangeStr)
	if matches == nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid range format11: %s", rangeStr)
	}
	startCol, startRow, err := excelize.CellNameToCoordinates(matches[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid range format66:%s,%s, %s", rangeStr, matches[1], err.Error())
	}
	if matches[2] == "" {
		return startCol, startRow, startCol, startRow, nil
	}

	endCol, endRow, err := excelize.CellNameToCoordinates(strings.TrimPrefix(matches[2], ":"))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid range format77:%s,%s, %s", rangeStr, matches[2], err.Error())
	}
	return startCol, startRow, endCol, endRow, nil
}
