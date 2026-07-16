package tools_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yz4230/hpc-wiki-mcp/pkg/tools"
	"github.com/yz4230/hpc-wiki-mcp/pkg/wiki"
)

func TestMCPServerSearchAndGetWikiPage(t *testing.T) {
	expectedAuthorization := "Basic " + base64.StdEncoding.EncodeToString([]byte("tester:secret"))
	wikiServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != expectedAuthorization {
			t.Errorf("Authorization = %q", got)
		}
		response.Header().Set("Content-Type", "text/html; charset=utf-8")
		if request.URL.Query().Get("do") == "search" {
			_, _ = response.Write([]byte(`<dl class="search_results"><div class="search_fullpage_result">
				<dt><a data-wiki-id="life.html" href="/doku.php?id=life.html">研究室生活Tips</a></dt>
				<dd class="meta"><span class="hits">3 ヒット</span></dd><dd class="snippet">施錠について</dd>
				</div></dl>`))
			return
		}
		_, _ = response.Write([]byte(`<h1>研究室生活Tips</h1><h2>退室</h2><p>最後の人は施錠。</p>`))
	}))
	t.Cleanup(wikiServer.Close)

	baseURL, _ := url.Parse(wikiServer.URL + "/doku.php")
	wikiClient := wiki.NewClient(&wiki.Config{BaseURL: baseURL, Username: "tester", Password: "secret", Timeout: time.Second}, wikiServer.Client())
	server := mcp.NewServer(&mcp.Implementation{Name: "hpc-wiki-mcp-server-test", Version: "0.1.0"}, nil)
	tools.Register(server, wikiClient)

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
	mcpServer := httptest.NewServer(handler)
	t.Cleanup(mcpServer.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "hpc-wiki-mcp-client-test", Version: "0.1.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: mcpServer.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	listed, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}
	names := map[string]bool{}
	for _, tool := range listed.Tools {
		names[tool.Name] = true
	}
	for _, name := range []string{"search_wiki", "get_wiki_page"} {
		if !names[name] {
			t.Fatalf("tool %q was not registered", name)
		}
	}

	searchResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "search_wiki", Arguments: map[string]any{"query": "施錠", "limit": 5},
	})
	if err != nil || searchResult.IsError {
		t.Fatalf("search_wiki = %#v, %v", searchResult, err)
	}
	var searchOutput tools.SearchWikiOutput
	decodeStructuredContent(t, searchResult.StructuredContent, &searchOutput)
	if len(searchOutput.Results) != 1 || searchOutput.Results[0].ID != "life.html" {
		t.Fatalf("search output = %#v", searchOutput)
	}
	if !strings.Contains(textContent(t, searchResult), "life.html") {
		t.Fatal("search text content does not include the page ID")
	}

	pageResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "get_wiki_page", Arguments: map[string]any{"id": "life.html"},
	})
	if err != nil || pageResult.IsError {
		t.Fatalf("get_wiki_page = %#v, %v", pageResult, err)
	}
	var pageOutput tools.GetWikiPageOutput
	decodeStructuredContent(t, pageResult.StructuredContent, &pageOutput)
	if pageOutput.Page == nil || !strings.Contains(pageOutput.Page.Text, "最後の人は施錠") {
		t.Fatalf("page output = %#v", pageOutput)
	}
}

func decodeStructuredContent(t *testing.T, content any, target any) {
	t.Helper()
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
}

func textContent(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("tool result returned no content")
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("tool result content type = %T", result.Content[0])
	}
	return text.Text
}
