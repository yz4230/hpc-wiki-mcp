package wiki

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestClientSearchAndGetPage(t *testing.T) {
	expectedAuthorization := "Basic " + base64.StdEncoding.EncodeToString([]byte("tester:secret"))
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
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
	t.Cleanup(server.Close)

	baseURL, _ := url.Parse(server.URL + "/doku.php")
	client := NewClient(&Config{BaseURL: baseURL, Username: "tester", Password: "secret", Timeout: time.Second}, server.Client())
	results, err := client.Search(context.Background(), "施錠", 5)
	if err != nil || len(results) != 1 || results[0].ID != "life.html" {
		t.Fatalf("Search() = %#v, %v", results, err)
	}
	page, err := client.GetPage(context.Background(), "life.html")
	if err != nil || page.Title != "研究室生活Tips" {
		t.Fatalf("GetPage() = %#v, %v", page, err)
	}
}
