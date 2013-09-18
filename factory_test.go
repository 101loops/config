// +build unit

package config

import (
	. "launchpad.net/gocheck"
	"strings"
)

var configTests = []struct {
	inputFiles  string
	shouldFail  bool
	errContains string
}{
	{
		inputFiles: "fixtures/shared.conf,fixtures/sample.conf",
	},
	{
		inputFiles:  "fixtures/sample.conf",
		shouldFail:  true,
		errContains: "invalid reference",
	}, {
		inputFiles:  "fixtures/empty.conf",
		shouldFail:  true,
		errContains: "missing",
	},
	{
		inputFiles:  "fixtures/invalid_syntax.conf",
		shouldFail:  true,
		errContains: "could not parse",
	},
}

// TESTS ==========================================================================================

func (s *S) TestLoadingConfig(c *C) {

	for i, configTest := range configTests {
		_, err := loadFromFiles("test", configTest.inputFiles)

		if err != nil {
			if !configTest.shouldFail {
				c.Fatalf("%d. Error parsing config %v: %v", i, configTest.inputFiles, err)
			} else {
				if !strings.Contains(err.Error(), configTest.errContains) {
					c.Fatalf("%d. Expected error containing '%v', got: %v", i, configTest.errContains, err)
				}
			}
		}
	}
}
