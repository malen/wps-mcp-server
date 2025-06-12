package mcps

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type WpsMCPServer struct {
	server *server.MCPServer
	// tools  []*server.ServerTool
}

func New() *WpsMCPServer {
	s := &WpsMCPServer{
		server: server.NewMCPServer("WPS Helper", "0.0.1"),
	}

	tool1 := mcp.NewTool("current_time", mcp.WithDescription("現在時刻を返します"))

	s.server.AddTool(tool1, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jst := time.FixedZone("Asiz/Tokyo", 9*60*60)
		now := time.Now().In(jst)
		message := fmt.Sprintf("現在の時刻：%s", now.Format("2006-01-02 15:04:05"))
		fmt.Println(message)
		return mcp.NewToolResultText(message), nil
	})

	AddWriteValuesToSheetTool(s.server)

	// if err := server.ServeStdio(mcpserver); err != nil {
	// 	fmt.Printf("サーバーエラー:%v \n", err)
	// }

	return s
}

func (wps *WpsMCPServer) Start() {
	// TODO: stdio
	// return server.ServeStdio(wps.server)

	// TODO: Streamable-http 通过http://localhost:8000/mcp 通讯
	httpServer := server.NewStreamableHTTPServer(wps.server)
	httpServer.Start(":8000")

	// TODO: SSE服务 通过http://localhost:8000/sse 通讯
	// httpServer := server.NewSSEServer(wps.server)
	// httpServer.Start(":8000")
}
