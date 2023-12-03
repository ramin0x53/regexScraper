package scraper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/storage"
)

type CompiledPattern struct {
	Regex   *regexp.Regexp
	Include []*regexp.Regexp
	Exclude []*regexp.Regexp
}

type ScrapperCfg struct {
	ScrapperId            int
	EntryPoint            string
	ThreadsCount          int
	Storage               storage.Storage
	AllowDomain           []string
	DisAllowDomain        []string
	RandomUserAgent       bool
	RandomMobileUserAgent bool
	UserAgent             string
	SearchKeywords        *[]CompiledPattern
	Result                chan<- SearchResult
}

func (cfg *ScrapperCfg) InitScrapper() {
	collector := colly.NewCollector(colly.Async())

	collector.SetStorage(cfg.Storage)
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: cfg.ThreadsCount})

	if len(cfg.AllowDomain) > 0 {
		collector.AllowedDomains = cfg.AllowDomain
	}
	if len(cfg.DisAllowDomain) > 0 {
		collector.DisallowedDomains = cfg.DisAllowDomain
	}
	if cfg.UserAgent != "" && !cfg.RandomUserAgent && !cfg.RandomMobileUserAgent {
		collector.UserAgent = cfg.UserAgent
	}
	if cfg.RandomUserAgent {
		extensions.RandomUserAgent(collector)
	}
	if cfg.RandomMobileUserAgent {
		extensions.RandomMobileUserAgent(collector)
	}

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	collector.OnResponse(func(r *colly.Response) {
		for _, searchKeyword := range *cfg.SearchKeywords {
			results := findRegex(&searchKeyword, &r.Body)
			for _, res := range results {
				fmt.Println(res)
				cfg.Result <- SearchResult{
					ScrapperId: cfg.ScrapperId,
					Url:        r.Request.URL.String(),
					Word:       res,
				}
			}
		}

		//check for none html pages
		if !strings.Contains(r.Headers.Get("Content-Type"), "html") {
			urlRegex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
			matchUrls := urlRegex.FindAllString(string(r.Body), -1)
			for _, link := range matchUrls {
				r.Request.Visit(link)
			}
		}
	})

	collector.Visit(cfg.EntryPoint)
}

func findRegex(regex *CompiledPattern, body *[]byte) []string {
	text := string(*body)
	match := regex.Regex.FindAllString(text, -1)
	if regex.Include != nil {
		for _, include := range regex.Include {
			if !include.MatchString(text) {
				return []string{}
			}
		}
	}
	if regex.Exclude != nil {
		for _, exclude := range regex.Exclude {
			if exclude.MatchString(text) {
				return []string{}
			}
		}
	}
	return match
}
