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

type SearchWikiInput struct {
	Query string `json:"query" jsonschema:"Japanese or English search query"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum results to return, from 1 to 50; defaults to 10"`
}

type SearchWikiOutput struct {
	Query   string              `json:"query"`
	Results []wiki.SearchResult `json:"results"`
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
