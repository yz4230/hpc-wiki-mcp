package wiki

type Link struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	External bool   `json:"external"`
}

type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

type Page struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Headings  []Heading `json:"headings"`
	Links     []Link    `json:"links"`
	SourceURL string    `json:"source_url"`
}

type SearchResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Snippet   string `json:"snippet"`
	Hits      int    `json:"hits,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	SourceURL string `json:"source_url"`
}
