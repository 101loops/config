package config

import (
	. "github.com/101loops/bdd"
	"fmt"
	"os"
)

var testConf *Config

func init() {
	os.Setenv("MEMCACHE_PASSWORD", "secretpassword")

	conf, err := loadFromFiles("test", "fixtures/shared.conf,fixtures/sample.conf")
	if err != nil {
		panic(fmt.Errorf("Error loading sample config: $v", err))
	}

	testConf = conf
}

var _ = Describe("Config", func() {

	It("content", func() {

		// ==== extract single section

			mail := testConf.SectionMust("mail")
			mailPort, _ := mail.Int("port")
			Check(mailPort, Equals, 25)
			mailHost, _ := mail.String("host")
			Check(mailHost, Equals, "smtp.mail.com")

			file := testConf.SectionMust("file")
			fileSystem, _ := file.String("system")
			Check(fileSystem, Equals, "s3")
			fileTmp, _ := file.Bool("tmp")
			Check(fileTmp, Equals, true)

			// ==== extract multiple sections

			caches := testConf.Sections("memcache")
			Check(caches, HasLen, 2)

			memcache := caches["default"]
			Check(memcache, NotNil)
			cacheHost, _ := memcache.String("host")
			Check(cacheHost, Equals, "memcache.com")
			cachePass, _ := memcache.String("pass")
			Check(cachePass, Equals, "secretpassword")

			memcache = caches["backup"]
			Check(memcache, NotNil)
			cacheHost, _ = memcache.String("host")
			//Check(cacheHost, Equals, "memcache.com")
			cachePass, _ = memcache.String("pass")
			//Check(cachePass, Equals, "secretpassword")

			// ==== extract missing section

			Check(func() { testConf.SectionMust("nonsense") }, Panics)
	})

	It("env", func() {
		cf, err := loadFromString("", "[profile] \n env: ")
		Check(err, IsNil)
		Check(cf.Env().IsDev(), Equals, true)
		Check(cf.Env().IsProd(), Equals, false)
		Check(cf.Env().IsTest(), Equals, false)
		Check(cf.Env().IsStage(), Equals, false)

		cf, err = loadFromString("", "[profile] \n env: staging")
		Check(err, IsNil)
		Check(cf.Env().IsDev(), Equals, false)
		Check(cf.Env().IsProd(), Equals, false)
		Check(cf.Env().IsTest(), Equals, false)
		Check(cf.Env().IsStage(), Equals, true)

		cf, err = loadFromString("", "[profile] \n env: testing")
		Check(err, IsNil)
		Check(cf.Env().IsDev(), Equals, false)
		Check(cf.Env().IsProd(), Equals, false)
		Check(cf.Env().IsTest(), Equals, true)
		Check(cf.Env().IsStage(), Equals, false)

		cf, err = loadFromString("", "[profile] \n env: development")
		Check(err, IsNil)
		Check(cf.Env().IsDev(), Equals, true)
		Check(cf.Env().IsProd(), Equals, false)
		Check(cf.Env().IsTest(), Equals, false)
		Check(cf.Env().IsStage(), Equals, false)

		cf, err = loadFromString("", "[profile] \n env: production")
		Check(err, IsNil)
		Check(cf.Env().IsDev(), Equals, false)
		Check(cf.Env().IsProd(), Equals, true)
		Check(cf.Env().IsTest(), Equals, false)
		Check(cf.Env().IsStage(), Equals, false)
	})
})
