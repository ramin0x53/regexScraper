package scraper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexScraper/config"
	"regexp"
)

func Run(cfg *config.Config, outputPath string) {
	storage := NewStorage(cfg.Redis)
	compiledPatterns := []CompiledPattern{}
	for _, pattern := range cfg.Search {
		compiledPattern := CompiledPattern{}
		compiledPattern.Regex = regexp.MustCompile(pattern.Regex)
		for _, inc := range pattern.Include {
			compiledPattern.Include = append(compiledPattern.Include, regexp.MustCompile(inc))
		}
		for _, exc := range pattern.Exclude {
			compiledPattern.Exclude = append(compiledPattern.Exclude, regexp.MustCompile(exc))
		}
		compiledPatterns = append(compiledPatterns, compiledPattern)
	}

	results := make(chan SearchResult, 100)
	for i, entryPoint := range cfg.EntryPoints {
		scrapperConfig := ScrapperCfg{
			ScrapperId:            i,
			EntryPoint:            entryPoint,
			ThreadsCount:          cfg.ThreadsCount,
			Storage:               storage,
			AllowDomain:           cfg.AllowDomains,
			DisAllowDomain:        cfg.DisAllowDomains,
			RandomUserAgent:       cfg.RandomUserAgent,
			RandomMobileUserAgent: cfg.RandomMobileUserAgent,
			UserAgent:             cfg.UserAgent,
			SearchKeywords:        &compiledPatterns,
			Result:                results,
		}
		go scrapperConfig.InitScrapper()
	}

	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer func() {
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer:", err)
			return
		}
	}()

	for {
		output := <-results
		if cfg.JsonOutput {
			o, err := json.Marshal(&output)
			if err != nil {
				log.Println(err)
			}
			_, err = fmt.Fprintln(writer, string(o))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			_, err = fmt.Fprintln(writer, output.Word)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
}
