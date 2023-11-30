package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RedisCfg struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

type RegexKeyword struct {
	Regex   string   `yaml:"regex"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type Config struct {
	EntryPoints           []string       `yaml:"entryPoints"`
	AllowDomains          []string       `yaml:"allowDomains"`
	DisAllowDomains       []string       `yaml:"disAllowDomains"`
	ThreadsCount          int            `yaml:"threadsCount"`
	Redis                 RedisCfg       `yaml:"redis"`
	JsonOutput            bool           `yaml:"jsonOutput"`
	RandomUserAgent       bool           `yaml:"randomUserAgent"`
	RandomMobileUserAgent bool           `yaml:"randomMobileUserAgent"`
	UserAgent             string         `yaml:"userAgent"`
	Search                []RegexKeyword `yaml:"search"`
}

func ReadConfig(path string) (*Config, error) {
	cfg := Config{}
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
