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

// PUBLIC METHODS =================================================================================

// Env returns the configuration's environment value (found in default section 'profile').
func (self *Config) Env() Env {
	env, _ := self.Sections(defaultSection)["default"].String("env")
	return Env(env)
}

// Version returns the configuration's version value (found in default section 'profile').
func (self *Config) Version() string {
	ver, _ := self.Sections(defaultSection)["default"].String("version")
	return ver
}

func (self *Config) Sections(want string) map[string]*Section {
	matches := make(map[string]*Section)
	for _, section := range self.fileConf.Sections() {
		if section == want || strings.HasPrefix(section, want+separator) {
			id := strings.Replace(strings.Replace(section, want+separator, "", 1), want, "", 1)
			if id == "" {
				id = "default"
			}
			matches[id] = &Section{
				id:   id,
				name: want,
				conf: self,
			}
		}
	}
	return matches
}

// HasSection tests whether a section exists.
func (self *Config) HasSection(want string) (match bool) {
	return len(self.Sections(want)) > 0
}

// SectionMust returns a section from the configuration;
// or panics if it does not exist.
func (self *Config) SectionMust(want string) (s *Section) {
	sections := self.Sections(want)
	if len(sections) == 0 {
		panic("config file does not have a section '" + want + "'")
	}
	if len(sections) > 1 {
		panic("config file has multiple sections '" + want + "'")
	}

	return sections["default"]
}
