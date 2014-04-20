package config

import (
	"fmt"
	. "github.com/101loops/bdd"
	"os"
)

var testConf *Config

func init() {
	os.Setenv("test.MEMCACHE_PASSWORD", "secretpassword")

	conf, err := loadFromFiles("test", "fixtures/shared.conf", "fixtures/sample.conf")
	if err != nil {
		panic(fmt.Errorf("error loading sample config: %v", err))
	}
	testConf = conf
}

var _ = Describe("Config", func() {

	With("extract options", func() {

		With("single section", func() {

			It("section 'mail'", func() {
				section := testConf.SectionMust("mail")

				port, err := section.Int("port")
				Check(err, IsNil)
				Check(port, Equals, 25)

				host, err := section.String("host")
				Check(err, IsNil)
				Check(host, Equals, "smtp.mail.com")

				//_, err = mail.String("pass")
				//Check(err, NotNil)
			})

			It("section 'file'", func() {
				section := testConf.SectionMust("file")

				system, err := section.String("system")
				Check(err, IsNil)
				Check(system, Equals, "s3")

				tmp, err := section.Bool("tmp")
				Check(err, IsNil)
				Check(tmp, IsTrue)
			})
		})

		It("multi-section", func() {
			sections := testConf.Sections("memcache")
			Check(sections, HasLen, 2)

			memcache := sections["default"]
			Check(memcache, NotNil)

			cacheHost, err := memcache.String("host")
			Check(err, IsNil)
			Check(cacheHost, Equals, "memcache.com")

			cachePass, err := memcache.String("pass")
			Check(err, IsNil)
			Check(cachePass, Equals, "secretpassword")

			memcache = sections["backup"]
			Check(memcache, NotNil)

			cacheHost, err = memcache.String("host")
			//Check(err, IsNil)
			//Check(cacheHost, Equals, "memcache.com")

			cachePass, err = memcache.String("pass")
			//Check(err, IsNil)
			//Check(cachePass, Equals, "secretpassword")
		})

		It("missing section", func() {
			Check(func() { testConf.SectionMust("nonsense") }, Panics)
		})
	})
})
