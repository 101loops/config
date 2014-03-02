package config

import (
	goconf "bitbucket.org/gosimple/conf"
	"strings"
)

const (
	separator      = ":"
	defaultSection = "profile"
)

type Config struct {
	fileConf *goconf.Config
	prefix   string
}

// Env returns the configuration's environment value (found in default section 'profile').
func (conf *Config) Env() Env {
	env, _ := conf.Sections(defaultSection)["default"].String("env")
	return Env(env)
}

// Version returns the configuration's version value (found in default section 'profile').
func (conf *Config) Version() string {
	ver, _ := conf.Sections(defaultSection)["default"].String("version")
	return ver
}

func (conf *Config) Sections(want string) map[string]*Section {
	matches := make(map[string]*Section)
	for _, section := range conf.fileConf.Sections() {
		if section == want || strings.HasPrefix(section, want+separator) {
			id := strings.Replace(strings.Replace(section, want+separator, "", 1), want, "", 1)
			if id == "" {
				id = "default"
			}
			matches[id] = &Section{
				id:   id,
				name: want,
				conf: conf,
			}
		}
	}
	return matches
}

// HasSection tests whether a section exists.
func (conf *Config) HasSection(want string) (match bool) {
	return len(conf.Sections(want)) > 0
}

// SectionMust returns a section from the configuration;
// or panics if it does not exist.
func (conf *Config) SectionMust(want string) (s *Section) {
	sections := conf.Sections(want)
	if len(sections) == 0 {
		panic("config file does not have a section '" + want + "'")
	}
	if len(sections) > 1 {
		panic("config file has multiple sections '" + want + "'")
	}

	return sections["default"]
}
