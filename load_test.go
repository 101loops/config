package config

import (
	. "github.com/101loops/bdd"
	"strings"
)

var configTests = []struct {
	inputFiles    []string
	shouldSucceed bool
	errContains   string
}{
	{
		inputFiles:    []string{"fixtures/shared.conf", "fixtures/sample.conf"},
		shouldSucceed: true,
	},
	{
		inputFiles:  []string{"fixtures/invalid_syntax.conf"},
		errContains: "could not parse line: [mail",
	},
	{
		inputFiles:  []string{"fixtures/invalid_syntax2.conf"},
		errContains: `invalid option name 'user.name': character '.' not allowed`,
	},
	{
		inputFiles:  []string{"fixtures/invalid_include.conf"},
		errContains: `invalid section include "mail.production": section not found`,
	},
	{
		inputFiles:  []string{"fixtures/invalid_reference.conf"},
		errContains: `section 'server' not found in "$server.port"`,
	},
}

var _ = Describe("Constructor", func() {

	Should("load config from file", func() {
		for i, configTest := range configTests {
			_, err := loadFromFiles("test", configTest.inputFiles...)

			if err != nil {
				if configTest.shouldSucceed {
					CheckFail("%d. Error parsing config %v: %v", i, configTest.inputFiles, err)
				} else {
					if !strings.Contains(err.Error(), configTest.errContains) {
						CheckFail("%d. Expected error to contain %q, got: %q", i, configTest.errContains, err)
					}
				}
			}
		}
	})
})
