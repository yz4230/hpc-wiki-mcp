package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yz4230/hpc-wiki-mcp/pkg/wiki"
)

const maxSearchContentItems = 20

var SearchWikiTool = &mcp.Tool{
	Name:        "search_wiki",
	Title:       "Search HPC Laboratory Wiki",
	Description: "Search the laboratory DokuWiki and return matching page IDs, titles, snippets, hit counts, update times, and source URLs.",
	Annotations: &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: new(false),
		IdempotentHint:  true,
		OpenWorldHint:   new(true),
	},
}

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

type SearchWikiInput struct {
	Query string `json:"query" jsonschema:"Japanese or English search query"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum results to return, from 1 to 50; defaults to 10"`
}

type SearchWikiOutput struct {
	Query   string              `json:"query"`
	Results []wiki.SearchResult `json:"results"`
}

type GetWikiPageInput struct {
	ID string `json:"id" jsonschema:"DokuWiki page ID, for example life.html or network:remote"`
}

type GetWikiPageOutput struct {
	Page *wiki.Page `json:"page"`
}

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

func formatSearchContent(query string, results []wiki.SearchResult) string {
	if len(results) == 0 {
		return fmt.Sprintf("No Wiki pages matched %q.", query)
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "Found %d Wiki pages matching %q.\n\n", len(results), query)
	limit := min(len(results), maxSearchContentItems)
	for index, result := range results[:limit] {
		fmt.Fprintf(&builder, "%d. [%s] %s", index+1, result.ID, result.Title)
		if result.Hits > 0 {
			fmt.Fprintf(&builder, " (%d hits)", result.Hits)
		}
		builder.WriteByte('\n')
		if result.Snippet != "" {
			fmt.Fprintf(&builder, "   %s\n", result.Snippet)
		}
	}
	if len(results) > limit {
		fmt.Fprintf(&builder, "\nShowing the first %d results. StructuredContent contains the full result.", limit)
	}
	return builder.String()
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

func (h *Handler) SearchWiki(ctx context.Context, _ *mcp.CallToolRequest, input *SearchWikiInput) (*mcp.CallToolResult, *SearchWikiOutput, error) {
	if input == nil {
		input = &SearchWikiInput{}
	}
	results, err := h.wiki.Search(ctx, input.Query, input.Limit)
	if err != nil {
		return nil, nil, err
	}
	return textResult(formatSearchContent(input.Query, results)), &SearchWikiOutput{
		Query:   input.Query,
		Results: results,
	}, nil
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
