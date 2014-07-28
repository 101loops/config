package config

import (
	"io/ioutil"

	goconf "github.com/gosimple/conf"
)

// New creates a Config from config files;
// or returns an error if any of the files couldn't be opened/parsed.
//
// The passed-in prefix is used to filter environment variables
// meant to overwrite config values of the files.
func New(prefix string, sources ...string) (*Config, error) {
	return loadFromFiles(prefix, sources...)
}

func loadFromFiles(prefix string, sources ...string) (*Config, error) {
	var confStr string
	for _, source := range sources {
		content, err := ioutil.ReadFile(source)
		if err != nil {
			return nil, err
		}
		confStr += string(content) + "\n"
	}
	return loadFromString(prefix, confStr)
}

func loadFromString(prefix string, confStr string) (*Config, error) {
	fileConf, err := goconf.ReadBytes([]byte(confStr))
	if err != nil {
		return nil, err
	}

	conf, err := newConf(prefix, fileConf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
