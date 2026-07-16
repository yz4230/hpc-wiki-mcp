package wiki

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	spaceRE    = regexp.MustCompile(`[ \t]+`)
	newlineRE  = regexp.MustCompile(`\n{3,}`)
	hitCountRE = regexp.MustCompile(`([\d,]+)\s*ヒット`)
)

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(spaceRE.ReplaceAllString(lines[i], " "))
	}
	return strings.TrimSpace(newlineRE.ReplaceAllString(strings.Join(lines, "\n"), "\n\n"))
}

func ParsePage(reader io.Reader, id string, sourceURL *url.URL) (*Page, error) {
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("parse Wiki page: %w", err)
	}
	document.Find("script, style, noscript").Remove()

	title := cleanText(document.Find("h1").First().Text())
	if title == "" {
		title = id
	}

	headings := make([]Heading, 0)
	document.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, selection *goquery.Selection) {
		text := cleanText(selection.Text())
		if text == "" {
			return
		}
		level, _ := strconv.Atoi(strings.TrimPrefix(goquery.NodeName(selection), "h"))
		headings = append(headings, Heading{Level: level, Text: text})
	})

	links := make([]Link, 0)
	seen := make(map[string]struct{})
	document.Find("a[href]").Each(func(_ int, selection *goquery.Selection) {
		rawHref, _ := selection.Attr("href")
		if rawHref == "" || strings.HasPrefix(rawHref, "#") {
			return
		}
		reference, err := url.Parse(rawHref)
		if err != nil {
			return
		}
		resolved := sourceURL.ResolveReference(reference)
		if resolved.Scheme != "http" && resolved.Scheme != "https" {
			return
		}
		title := cleanText(selection.Text())
		if title == "" {
			title = resolved.String()
		}
		key := resolved.String() + "\x00" + title
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		wikiID, _ := selection.Attr("data-wiki-id")
		if wikiID == "" {
			wikiID = reference.Query().Get("id")
		}
		links = append(links, Link{
			ID:       wikiID,
			Title:    title,
			URL:      resolved.String(),
			External: resolved.Scheme != sourceURL.Scheme || resolved.Host != sourceURL.Host,
		})
	})

	return &Page{
		ID:        id,
		Title:     title,
		Text:      cleanText(document.Text()),
		Headings:  headings,
		Links:     links,
		SourceURL: sourceURL.String(),
	}, nil
}

func ParseSearch(reader io.Reader, baseURL *url.URL, limit int) ([]SearchResult, error) {
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("parse Wiki search results: %w", err)
	}

	results := make([]SearchResult, 0, limit)
	document.Find(".search_results .search_fullpage_result").EachWithBreak(func(_ int, result *goquery.Selection) bool {
		if len(results) >= limit {
			return false
		}
		link := result.Find("dt a[data-wiki-id]").First()
		id, _ := link.Attr("data-wiki-id")
		href, _ := link.Attr("href")
		if id == "" || href == "" {
			return true
		}
		reference, err := url.Parse(href)
		if err != nil {
			return true
		}
		hits := 0
		if match := hitCountRE.FindStringSubmatch(cleanText(result.Find("dd.meta").Text())); len(match) == 2 {
			hits, _ = strconv.Atoi(strings.ReplaceAll(match[1], ",", ""))
		}
		updatedAt, _ := result.Find("dd.meta time").Attr("datetime")
		results = append(results, SearchResult{
			ID:        id,
			Title:     cleanText(link.Text()),
			Snippet:   cleanText(result.Find("dd.snippet").Text()),
			Hits:      hits,
			UpdatedAt: updatedAt,
			SourceURL: baseURL.ResolveReference(reference).String(),
		})
		return true
	})
	return results, nil
}
