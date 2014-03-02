// +build unit

package config

import (
	. "launchpad.net/gocheck"
	"os"
)

var testConf *Config

func (s *S) SetUpTest(c *C) {
	os.Setenv("MEMCACHE_PASSWORD", "secretpassword")

	conf, err := loadFromFiles("test", "fixtures/shared.conf,fixtures/sample.conf")
	if err != nil {
		c.Fatalf("Error loading sample config: $v", err)
	}

	testConf = conf
}

func (s *S) TestConfigContent(c *C) {

	// ==== extract single section

	mail := testConf.SectionMust("mail")
	mailPort, _ := mail.Int("port")
	c.Assert(mailPort, Equals, 25)
	mailHost, _ := mail.String("host")
	c.Assert(mailHost, Equals, "smtp.mail.com")

	file := testConf.SectionMust("file")
	fileSystem, _ := file.String("system")
	c.Assert(fileSystem, Equals, "s3")
	fileTmp, _ := file.Bool("tmp")
	c.Assert(fileTmp, Equals, true)

	// ==== extract multiple sections

	caches := testConf.Sections("memcache")
	c.Assert(caches, HasLen, 2)

	memcache := caches["default"]
	c.Assert(memcache, NotNil)
	cacheHost, _ := memcache.String("host")
	c.Assert(cacheHost, Equals, "memcache.com")
	cachePass, _ := memcache.String("pass")
	c.Assert(cachePass, Equals, "secretpassword")

	memcache = caches["backup"]
	c.Assert(memcache, NotNil)
	cacheHost, _ = memcache.String("host")
	//c.Assert(cacheHost, Equals, "memcache.com")
	cachePass, _ = memcache.String("pass")
	//c.Assert(cachePass, Equals, "secretpassword")

	// ==== extract missing section

	c.Assert(func() { testConf.SectionMust("nonsense") }, PanicMatches, "*does not have a section 'nonsense'")
}

func (s *S) TestEnv(c *C) {
	cf, err := loadFromString("", "[profile] \n env: ")
	c.Assert(err, IsNil)
	c.Assert(cf.Env().IsDev(), Equals, true)
	c.Assert(cf.Env().IsProd(), Equals, false)
	c.Assert(cf.Env().IsTest(), Equals, false)
	c.Assert(cf.Env().IsStage(), Equals, false)

	cf, err = loadFromString("", "[profile] \n env: staging")
	c.Assert(err, IsNil)
	c.Assert(cf.Env().IsDev(), Equals, false)
	c.Assert(cf.Env().IsProd(), Equals, false)
	c.Assert(cf.Env().IsTest(), Equals, false)
	c.Assert(cf.Env().IsStage(), Equals, true)

	cf, err = loadFromString("", "[profile] \n env: testing")
	c.Assert(err, IsNil)
	c.Assert(cf.Env().IsDev(), Equals, false)
	c.Assert(cf.Env().IsProd(), Equals, false)
	c.Assert(cf.Env().IsTest(), Equals, true)
	c.Assert(cf.Env().IsStage(), Equals, false)

	cf, err = loadFromString("", "[profile] \n env: development")
	c.Assert(err, IsNil)
	c.Assert(cf.Env().IsDev(), Equals, true)
	c.Assert(cf.Env().IsProd(), Equals, false)
	c.Assert(cf.Env().IsTest(), Equals, false)
	c.Assert(cf.Env().IsStage(), Equals, false)

	cf, err = loadFromString("", "[profile] \n env: production")
	c.Assert(err, IsNil)
	c.Assert(cf.Env().IsDev(), Equals, false)
	c.Assert(cf.Env().IsProd(), Equals, true)
	c.Assert(cf.Env().IsTest(), Equals, false)
	c.Assert(cf.Env().IsStage(), Equals, false)
}
