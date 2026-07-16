package wiki

import (
	"net/url"
	"strings"
	"testing"
)

func TestNormalizePageID(t *testing.T) {
	for input, want := range map[string]string{
		":network:remote:": "network:remote",
		"life.html":        "life.html",
	} {
		got, err := NormalizePageID(input)
		if err != nil || got != want {
			t.Fatalf("NormalizePageID(%q) = %q, %v; want %q", input, got, err, want)
		}
	}
	for _, input := range []string{"", "../secret", "network/remote"} {
		if _, err := NormalizePageID(input); err == nil {
			t.Fatalf("NormalizePageID(%q) succeeded", input)
		}
	}
}

func TestParsePage(t *testing.T) {
	sourceURL, _ := url.Parse("http://wiki.example/wiki/doku.php?id=life.html&do=export_xhtmlbody")
	page, err := ParsePage(strings.NewReader(`<h1>生活 Tips</h1><h2>入室</h2><p>最後の人は施錠。</p>
		<a data-wiki-id="network:remote" href="/wiki/doku.php?id=network:remote">リモート</a>
		<a href="https://example.com/guide">外部資料</a>`), "life.html", sourceURL)
	if err != nil {
		t.Fatalf("ParsePage() error = %v", err)
	}
	if page.Title != "生活 Tips" || !strings.Contains(page.Text, "最後の人は施錠") {
		t.Fatalf("unexpected page: %#v", page)
	}
	if len(page.Headings) != 2 || page.Headings[1].Level != 2 {
		t.Fatalf("headings = %#v", page.Headings)
	}
	if len(page.Links) != 2 || page.Links[0].ID != "network:remote" || page.Links[0].External || !page.Links[1].External {
		t.Fatalf("links = %#v", page.Links)
	}
}

func TestParseSearch(t *testing.T) {
	baseURL, _ := url.Parse("http://wiki.example/wiki/doku.php")
	results, err := ParseSearch(strings.NewReader(`<dl class="search_results">
		<div class="search_fullpage_result">
		<dt><a data-wiki-id="research.html" href="/wiki/doku.php?id=research.html">研究Tips</a></dt>
		<dd class="meta"><span class="hits">16 ヒット</span><time datetime="2026-01-02T03:04:05+00:00">更新</time></dd>
		<dd class="snippet">論文を <strong>探す</strong> 方法</dd></div></dl>`), baseURL, 10)
	if err != nil {
		t.Fatalf("ParseSearch() error = %v", err)
	}
	if len(results) != 1 || results[0].ID != "research.html" || results[0].Hits != 16 || results[0].Snippet != "論文を 探す 方法" {
		t.Fatalf("results = %#v", results)
	}
}
