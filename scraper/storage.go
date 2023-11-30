package scraper

import (
	"regexScraper/config"

	"github.com/gocolly/colly/v2/storage"
	"github.com/gocolly/redisstorage"
)

func NewStorage(redisCfg config.RedisCfg) storage.Storage {
	if redisCfg.Address != "" {
		return &redisstorage.Storage{
			Address:  redisCfg.Address,
			Password: redisCfg.Password,
			DB:       redisCfg.Db,
			Prefix:   redisCfg.Prefix,
		}
	}

	storage := storage.InMemoryStorage{}
	storage.Init()
	return &storage
}
