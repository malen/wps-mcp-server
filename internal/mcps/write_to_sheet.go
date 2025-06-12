package mcps

import (
	"aoisoft/wps-mcp-server/internal/utils"
	"context"
	"fmt"

	z "github.com/Oudwins/zog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/xuri/excelize/v2"
)

type writeValuesToSheetArgs struct {
	FileAbsolutePath string     `zog:"fileAbsolutePath"`
	SheetName        string     `zog:"sheetName"`
	NewSheet         bool       `zog:"newSheet"`
	Range            string     `zog:"range"`
	Values           [][]string `zog:"values"`
}

var writeValuesToSheetSchema = z.Struct(z.Shape{
	"fileAbsolutePath": z.String().Test(utils.AbsolutePathTest()).Required(),
	"sheetName":        z.String().Required(),
	"newSheet":         z.Bool().Required().Default(false),
	"range":            z.String().Required(),
	"values":           z.Slice(z.Slice(z.String())).Required(),
})

func AddWriteValuesToSheetTool(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("write_values_to_sheet",
		mcp.WithDescription("Write values to the sheet of Wps"),
		mcp.WithString("fileAbsolutePath",
			mcp.Required(), mcp.Description("The absolute path to the WPS file"),
		),
		mcp.WithString("sheetName",
			mcp.Required(),
			mcp.Description("Sheet name in the WPS file"),
		),
		mcp.WithBoolean("newSheet",
			mcp.Required(),
			mcp.Description("create a new sheet if true, otherwise write values to a existing sheet"),
		),
		mcp.WithString("range",
			mcp.Required(),
			mcp.Description("range of cells in the WPS sheet (e.g., \"A1:B2\")"),
		),
		mcp.WithArray("values",
			mcp.Required(),
			mcp.Description("values to write to the WPS sheet. If the value is a forula, it should start with \"=\""),
			mcp.Items(map[string]any{
				"type": "array",
				"items": map[string]any{
					"anyOf": []any{
						map[string]any{
							"type": "string",
						},
						map[string]any{
							"type": "number",
						},
						map[string]any{
							"type": "boolean",
						},

						map[string]any{
							"type": "null",
						},
					},
				},
			}),
		),
	), writeValuesToSheetHandler)
}

func writeValuesToSheetHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := writeValuesToSheetArgs{}
	issues := writeValuesToSheetSchema.Parse(request.GetArguments(), &args)
	if len(issues) != 0 {
		return NewToolResultZogIssueMap(issues), nil
	}

	wps, closeFunc, err := OpenFile(args.FileAbsolutePath)
	if err != nil {
		return nil, err
	}
	defer closeFunc()

	if args.NewSheet {
		if err := wps.AddNewSheet(args.SheetName); err != nil {
			return nil, err
		}
	}

	sheet, err := wps.FindSheet(args.SheetName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid argument11: %s", err.Error())), nil
	}

	startCol, startRow, endCol, endRow, err := utils.ParseRange(args.Range)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid argument22: %s %s", args.Range, err.Error())), nil
	}

	// データの整合性チェック
	rangeRowSize := endRow - startRow + 1
	if len(args.Values) != rangeRowSize {
		return mcp.NewToolResultError(fmt.Sprintf("number of rows in data (%d) does not match range size (%d)", len(args.Values), rangeRowSize)), nil
	}

	for i, row := range args.Values {
		rangeColSize := endCol - startCol + 1
		if len(row) != rangeColSize {
			return nil, nil
		}
		for j, cellValue := range row {
			cell, err := excelize.CoordinatesToCellName(startCol+j, startRow+i)
			if err != nil {
				return nil, err
			}
			fmt.Printf("xxxx cell value : %s, %s:\n", cell, cellValue)
			err = sheet.SetValue(cell, cellValue)

			if err != nil {
				return nil, err
			}
		}

	}

	if err := wps.Save(); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText("success"), nil
}
