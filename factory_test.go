package config

import (
	. "github.com/101loops/bdd"
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

var _ = Describe("Factory", func() {

	It("create config from file", func() {
		for i, configTest := range configTests {
			_, err := loadFromFiles("test", configTest.inputFiles)

			if err != nil {
				if !configTest.shouldFail {
					CheckFail("%d. Error parsing config %v: %v", i, configTest.inputFiles, err)
				} else {
					if !strings.Contains(err.Error(), configTest.errContains) {
						CheckFail("%d. Expected error containing %q, got: %v", i, configTest.errContains, err)
					}
				}
			}
		}
	})
})
