package scraper

type SearchResult struct {
	ScrapperId int    `json:"scrapperId"`
	Url        string `json:"url"`
	Word       string `json:"word"`
}
