package mcps

import (
	"fmt"

	z "github.com/Oudwins/zog"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewToolResultZogIssueMap(errs z.ZogIssueMap) *mcp.CallToolResult {
	issues := z.Issues.SanitizeMap(errs)

	var issueResults []mcp.Content
	for k, messages := range issues {
		for _, message := range messages {
			issueResults = append(issueResults, mcp.NewTextContent(fmt.Sprintf("Invalid argument33: %s: %s", k, message)))
		}
	}

	return &mcp.CallToolResult{
		Content: issueResults,
		IsError: true,
	}
}
