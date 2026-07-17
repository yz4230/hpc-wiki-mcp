package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yz4230/hpc-wiki-mcp/pkg/wiki"
)

var GetWikiPageTool = &mcp.Tool{
	Name:        "get_wiki_page",
	Title:       "Read an HPC Laboratory Wiki Page",
	Description: "Fetch a DokuWiki page by ID and return its text, headings, links, and source URL. Use search_wiki first when the page ID is unknown.",
	Annotations: &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: new(false),
		IdempotentHint:  true,
		OpenWorldHint:   new(true),
	},
}

type GetWikiPageInput struct {
	ID string `json:"id" jsonschema:"DokuWiki page ID, for example life.html or network:remote"`
}

type GetWikiPageOutput struct {
	Page *wiki.Page `json:"page"`
}

func formatPageContent(page *wiki.Page) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# %s\n\n", page.Title)
	fmt.Fprintf(&builder, "page_id: %s\nsource_url: %s\n", page.ID, page.SourceURL)
	if page.Text != "" {
		fmt.Fprintf(&builder, "\n%s", page.Text)
	}
	return builder.String()
}

func (h *Handler) GetWikiPage(ctx context.Context, _ *mcp.CallToolRequest, input *GetWikiPageInput) (*mcp.CallToolResult, *GetWikiPageOutput, error) {
	if input == nil {
		input = &GetWikiPageInput{}
	}
	page, err := h.wiki.GetPage(ctx, input.ID)
	if err != nil {
		return nil, nil, err
	}
	return textResult(formatPageContent(page)), &GetWikiPageOutput{Page: page}, nil
}
