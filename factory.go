package config

import (
	"log"
	goconf "bitbucket.org/gosimple/conf"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// FACTORY ========================================================================================

// NewConf loads one/multiple configuration files (paths separated by comma);
// or an error if any of the files couldn't be opened/read.
// A prefix allows to distinguish environment variables / command-line args meant
// to overwrite config values.
func NewConf(prefix string, defaultFile string) (*Config, error) {
	return load(prefix, defaultFile)
}

// HELPERS ========================================================================================

func load(prefix string, defaultFile string) (*Config, error) {
	configPath := flag.String("config", defaultFile, "path to the config file(s)")
	log.Printf("loading configuration file(s): '%v'", *configPath)
	return loadFromFiles(prefix, *configPath)
}

func loadFromFiles(prefix string, fileNames string) (*Config, error) {
	var configData string
	for _, fname := range strings.Split(fileNames, ",") {
		content, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, err
		}
		configData += string(content) + "\n"
	}
	return loadFromString(prefix, configData)
}

func loadFromString(prefix string, configData string) (*Config, error) {

	// create
	fileConf, err := goconf.ReadBytes([]byte(configData))
	if err != nil {
		return nil, err
	}

	// evaluate
	fileConf, err = evaluate(fileConf)
	if err != nil {
		return nil, err
	}

	// validate
	conf := &Config{fileConf: fileConf, prefix: prefix}
	err = validate(conf)

	return conf, err
}

func validate(c *Config) error {

	if !c.HasSection(defaultSection) {
		return fmt.Errorf("missing section '%s'", defaultSection)
	}

	env := c.Env()
	if env != "" && env != "development" && env != "production" && env != "staging" && env != "testing" {
		return fmt.Errorf("invalid application environment: %s", env)
	}

	return nil
}

func evaluate(conf *goconf.Config) (*goconf.Config, error) {
	for _, s := range conf.Sections() {
		opts, _ := conf.Options(s)
		for _, o := range opts {
			if conf.HasOption(s, o) {
				val, _ := conf.RawString(s, o)

				if o == "$" {
					// substitute section reference
					if conf.HasSection(val) {
						refOpts, _ := conf.Options(val)
						for _, refOpt := range refOpts {
							if conf.HasOption(val, refOpt) {
								refVal, _ := conf.RawString(val, refOpt)
								if !conf.HasOption(s, refOpt) {
									conf.AddOption(s, refOpt, refVal)
								}
							}
						}
					}
					conf.RemoveOption(s, "$")
				} else {
					// substitute option reference
					if strings.HasPrefix(val, "$") {
						ref, err := find(conf, val[1:])
						if err != nil {
							return conf, err
						}
						conf.AddOption(s, o, ref)
					}
				}
			}
		}
	}
	//print(conf)
	return conf, nil
}

func find(conf *goconf.Config, val string) (ref string, err error) {

	// lookup file
	parts := strings.Split(val, ".")
	if len(parts) > 1 {
		secRef := strings.Join(parts[:len(parts)-1], ".")
		optRef := parts[len(parts)-1]
		ref, _ = conf.RawString(secRef, optRef)
	}

	// lookup ENV
	if ref == "" {
		ref = os.Getenv(val)
	}

	if ref == "" {
		err = fmt.Errorf("invalid reference '%s'", val)
	}

	return
}

func print(conf *goconf.Config) {
	var b bytes.Buffer
	conf.Write(&b, "")
	b.WriteTo(os.Stdout)
}
