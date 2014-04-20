package config

import (
	goconf "github.com/gosimple/conf"
	"strconv"
	"strings"
	"time"
)

// Section represents a named group of configuration options.
type Section struct {
	id   string // e.g. "s3"
	name string // e.g. "files"
	conf *Config
}

// ID returns the unique identifier of the section ("default" by default).
func (sect *Section) ID() string {
	return sect.id
}

// Add adds a new option and value to the configuration.
// It returns true if the option and value were inserted,
// and false if the value was overwritten.
func (sect *Section) Add(option string, value string) bool {
	return sect.conf.fileConf.AddOption(sect.name, option, value)
}

// Bool returns the boolean value of the option;
// or an error if the conversion failed / the value didn't exist.
func (sect *Section) Bool(option string) (bool, error) {
	str, err := sect.get(option)
	if err != nil {
		return false, err
	}

	value, ok := goconf.BoolStrings[strings.ToLower(str)]
	if !ok {
		return false, goconf.GetError{goconf.CouldNotParse, "bool", str, sect.name, option}
	}

	return value, nil
}

// Float returns the decimal value of the option;
// or an error if the conversion failed / the value didn't exist.
func (sect *Section) Float(option string) (float64, error) {
	str, err := sect.get(option)
	if err != nil {
		return 0.0, err
	}

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		err = goconf.GetError{goconf.CouldNotParse, "float", str, sect.name, option}
	}
	return value, err
}

// Int returns the numeric value of the option;
// or an error if the conversion failed / the value didn't exist.
func (sect *Section) Int(option string) (int, error) {
	str, err := sect.get(option)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseInt(str, 10, 32)
	return int(parsed), err
}

// String returns the text value of the option;
// or an error if it does not exist.
func (sect *Section) String(option string) (string, error) {
	return sect.get(option)
}

// Bytes returns the raw byte array of the option;
// or an error if the conversion failed / the value didn't exist.
func (sect *Section) Bytes(option string) ([]byte, error) {
	str, err := sect.get(option)
	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}

// Duration returns the time duration value of the option;
// or an error if the conversion failed / the value didn't exist.
func (sect *Section) Duration(option string) (time.Duration, error) {
	str, err := sect.get(option)
	if err != nil {
		return time.Duration(0), err
	}
	return time.ParseDuration(str)
}

func (sect *Section) get(option string) (string, error) {
	name := sect.name
	if sect.id != "" && sect.id != "default" {
		name += "-" + sect.id
	}
	return sect.conf.get(name, option)
}
