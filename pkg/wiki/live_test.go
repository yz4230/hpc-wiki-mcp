package wiki

import (
	"context"
	"os"
	"testing"
)

func TestLiveWiki(t *testing.T) {
	if os.Getenv("HPC_WIKI_LIVE") != "1" {
		t.Skip("set HPC_WIKI_LIVE=1 to test the real Wiki")
	}
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	client := NewClient(config, nil)

	results, err := client.Search(context.Background(), "研究", 2)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) == 0 || results[0].ID == "" {
		t.Fatalf("Search() returned no usable results: %#v", results)
	}

	page, err := client.GetPage(context.Background(), "life.html")
	if err != nil {
		t.Fatalf("GetPage() error = %v", err)
	}
	if page.ID != "life.html" || page.Title == "" || page.Text == "" {
		t.Fatalf("GetPage() returned an incomplete page: %#v", page)
	}
}
