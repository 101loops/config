package config

import (
	"fmt"
	"os"
	"strings"

	goconf "github.com/gosimple/conf"
)

const (
	defaultSection        = "default"
	optionReferenceSymbol = "$"
	sectionIDSeparator    = ":"
	sectionImportKeyword  = "$include"
)

// Config represents configuration settings.
type Config struct {
	fileConf *goconf.Config
	prefix   string
}

func newConf(prefix string, fileConf *goconf.Config) (*Config, error) {
	conf := &Config{fileConf, prefix}

	for _, section := range fileConf.Sections() {
		opts, _ := fileConf.Options(section)
		for _, opt := range opts {
			val, _ := fileConf.RawString(section, opt)

			if strings.Contains(opt, ".") {
				return nil, fmt.Errorf("invalid option name '" + opt + "': character '.' not allowed")
			}

			if opt == sectionImportKeyword {
				if !fileConf.HasSection(val) {
					return nil, fmt.Errorf("invalid section include %q: section not found", val)
				}

				refOpts, _ := fileConf.Options(val)
				for _, refOpt := range refOpts {
					refVal, _ := fileConf.RawString(val, refOpt)
					if !fileConf.HasOption(section, refOpt) {
						fileConf.AddOption(section, refOpt, refVal)
					}
				}

				fileConf.RemoveOption(section, sectionImportKeyword)
				continue
			}

			if strings.HasPrefix(val, optionReferenceSymbol) {
				ref, err := conf.getPath(val[1:])
				if err != nil {
					return nil, fmt.Errorf("%s in %q", err, val)
				}
				fileConf.AddOption(section, opt, ref)
			}
		}
	}

	return conf, nil
}

// Sections returns a map of all sections with the passed-in name.
// The key of the returned map equals the section's ID.
func (conf *Config) Sections(name string) map[string]*Section {
	keyVal := make(map[string]*Section)
	for _, sectionName := range conf.fileConf.Sections() {
		if sectionName == name {
			keyVal[defaultSection] = &Section{
				id: defaultSection, name: name, conf: conf,
			}
		} else if strings.HasPrefix(sectionName, name+sectionIDSeparator) {
			id := strings.Replace(strings.Replace(sectionName, name+sectionIDSeparator, "", 1), name, "", 1)
			if id == "" {
				id = defaultSection
			}
			keyVal[id] = &Section{
				id: id, name: name, conf: conf,
			}
		}
	}
	return keyVal
}

// HasSection tests whether a section exists.
func (conf *Config) HasSection(name string) (match bool) {
	return len(conf.Sections(name)) > 0
}

// SectionMust returns the section from the config that matches the passed-in name;
// or panics if it does not exist. If there are multiple sections with the same name
// the optional passed-in id decides which to use.
func (conf *Config) SectionMust(name string, id ...string) *Section {
	sections := conf.Sections(name)

	if len(sections) == 0 {
		panic("config does not have a section '" + name + "'")
	}

	if len(sections) > 1 {
		if len(id) == 0 {
			panic("config has multiple sections '" + name + "'; an ID is required!")
		}
		sectionID := id[0]
		if s, ok := sections[sectionID]; ok {
			return s
		}
		panic("config has no section '" + name + "' with ID '" + sectionID + "'")
	}

	return sections[defaultSection]
}

func (conf *Config) getPath(path string) (string, error) {
	pathArr := strings.Split(path, ".")
	if len(pathArr) > 1 {
		section := strings.Join(pathArr[:len(pathArr)-1], ".")
		option := pathArr[len(pathArr)-1]
		return conf.get(section, option)
	}

	return conf.get(defaultSection, path)
}

func (conf *Config) get(section, option string) (string, error) {
	path := section + "." + option
	if section == "" || section == defaultSection {
		path = option
	}

	// 1) lookup command line arguments
	if val, ok := cmdArgs[path]; ok {
		return val, nil
	}

	// 2) lookup environment variable
	if val := os.Getenv(conf.prefix + "." + path); val != "" {
		return val, nil
	}

	// 3) lookup configuration file
	return conf.fileConf.String(section, option)
}
