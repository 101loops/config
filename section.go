package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Section struct {
	name string // e.g. files
	id   string // e.g. s3
	conf *Config
}

var (
	cmdArgs map[string]string
)

// Id returns the unique identifier of the section ("default" by default).
func (self *Section) Id() string {
	return self.id
}

// Add adds a new option and value to the configuration.
// It returns true if the option and value were inserted, and false if the value was overwritten.
func (self *Section) Add(option string, value string) bool {
	return self.conf.fileConf.AddOption(self.name, option, value)
}

// Bool returns the boolean value of the option;
// or an error if the conversion failed / the value didn't exist.
func (self *Section) Bool(option string) (res bool, err error) {
	str, err := self.get(option)
	if err != nil {
		return res, err
	}
	if str == "true" {
		return true, err
	}

	return false, err
}

// Float returns the decimal value of the option;
// or an error if the conversion failed / the value didn't exist.
func (self *Section) Float(option string) (res float64, err error) {
	str, err := self.get(option)
	if err == nil {
		res, err = strconv.ParseFloat(str, 64)
	}
	return res, err
}

// Int returns the numeric value of the option;
// or an error if the conversion failed / the value didn't exist.
func (self *Section) Int(option string) (res int, err error) {
	str, err := self.get(option)
	if err == nil {
		parsed, err := strconv.ParseInt(str, 10, 32)
		return int(parsed), err
	}
	return res, err
}

// String returns the text value of the option, or an error if it does not exist.
func (self *Section) String(option string) (string, error) {
	return self.get(option)
}

// Bytes returns the raw byte array of the option;
// or an error if the conversion failed / the value didn't exist.
func (self *Section) Bytes(option string) ([]byte, error) {
	s, err := self.get(option)
	if err == nil {
		return []byte(s), nil
	}
	return nil, err
}

// Duration returns the time duration value of the option;
// or an error if the conversion failed / the value didn't exist.
func (self *Section) Duration(option string) (d time.Duration, err error) {
	s, err := self.get(option)
	if err == nil {
		d, err = time.ParseDuration(s)
	}
	return
}

// PRIVATE METHODS ================================================================================

func (self *Section) get(option string) (string, error) {

	section := self.sectionName()
	path := toPath(section, option)

	// first: search command line arguments
	res, ok := getArgs()[path]
	if ok {
		return res, nil
	}

	// second: search environment variable
	res = os.Getenv(self.conf.prefix + "." + path)
	if res != "" {
		return res, nil
	}

	// finally: search file configuration
	return self.conf.fileConf.String(section, option)
}

func (self *Section) sectionName() (s string) {
	s = self.name
	if self.id != "" && self.id != "default" {
		s += "-" + self.id
	}
	return
}

func getArgs() map[string]string {
	if cmdArgs == nil {
		cmdArgs = make(map[string]string)
		var name string
		for _, arg := range os.Args {
			if strings.HasPrefix(arg, "-") {
				name = strings.TrimLeft(arg, "-")
			}
			if name != "" {
				cmdArgs[name] = arg
				name = ""
			}
		}
	}
	return cmdArgs
}

func toPath(section string, option string) string {
	if section == "" || section == defaultSection {
		return option
	}
	return section + "." + option
}
