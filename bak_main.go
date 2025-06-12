package bakmain

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpserver := server.NewMCPServer("WPS Helper", "0.0.1")

	tool1 := mcp.NewTool("current_time", mcp.WithDescription("現在時刻を返します"))

	mcpserver.AddTool(tool1, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jst := time.FixedZone("Asiz/Tokyo", 9*60*60)
		now := time.Now().In(jst)
		message := fmt.Sprintf("現在の時刻：%s", now.Format("2006-01-02 15:04:05"))
		return mcp.NewToolResultText(message), nil
	})

	if err := server.ServeStdio(mcpserver); err != nil {
		fmt.Printf("サーバーエラー:%v \n", err)
	}

}
