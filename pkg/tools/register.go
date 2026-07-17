package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yz4230/hpc-wiki-mcp/pkg/wiki"
)

type Handler struct {
	wiki *wiki.Client
}

func Register(server *mcp.Server, client *wiki.Client) {
	handler := &Handler{wiki: client}
	mcp.AddTool(server, SearchWikiTool, handler.SearchWiki)
	mcp.AddTool(server, GetWikiPageTool, handler.GetWikiPage)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}
