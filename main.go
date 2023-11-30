package main

import (
	"flag"
	"log"
	"regexScraper/config"
	"regexScraper/scraper"
)

func main() {
	configPath := flag.String("c", "", "configPath")
	outputPath := flag.String("o", "", "outputPath")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	scraper.Run(cfg, *outputPath)
}
