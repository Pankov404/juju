// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv_test

import (
	gc "launchpad.net/gocheck"

	"launchpad.net/juju-core/juju/osenv"
	"launchpad.net/juju-core/testing/testbase"
)

type proxySuite struct {
	testbase.LoggingSuite
}

var _ = gc.Suite(&proxySuite{})

func (s *proxySuite) TestDetectNoSettings(c *gc.C) {
	// Patch all of the environment variables we check out just in case the
	// user has one set.
	s.PatchEnvironment("http_proxy", "")
	s.PatchEnvironment("HTTP_PROXY", "")
	s.PatchEnvironment("https_proxy", "")
	s.PatchEnvironment("HTTPS_PROXY", "")
	s.PatchEnvironment("ftp_proxy", "")
	s.PatchEnvironment("FTP_PROXY", "")

	proxies := osenv.DetectProxies()

	c.Assert(proxies, gc.DeepEquals, osenv.ProxySettings{})
}

func (s *proxySuite) TestDetectPrimary(c *gc.C) {
	// Patch all of the environment variables we check out just in case the
	// user has one set.
	s.PatchEnvironment("http_proxy", "http://user@10.0.0.1")
	s.PatchEnvironment("HTTP_PROXY", "")
	s.PatchEnvironment("https_proxy", "https://user@10.0.0.1")
	s.PatchEnvironment("HTTPS_PROXY", "")
	s.PatchEnvironment("ftp_proxy", "ftp://user@10.0.0.1")
	s.PatchEnvironment("FTP_PROXY", "")

	proxies := osenv.DetectProxies()

	c.Assert(proxies, gc.DeepEquals, osenv.ProxySettings{
		Http:  "http://user@10.0.0.1",
		Https: "https://user@10.0.0.1",
		Ftp:   "ftp://user@10.0.0.1",
	})
}

func (s *proxySuite) TestDetectFallback(c *gc.C) {
	// Patch all of the environment variables we check out just in case the
	// user has one set.
	s.PatchEnvironment("http_proxy", "")
	s.PatchEnvironment("HTTP_PROXY", "http://user@10.0.0.2")
	s.PatchEnvironment("https_proxy", "")
	s.PatchEnvironment("HTTPS_PROXY", "https://user@10.0.0.2")
	s.PatchEnvironment("ftp_proxy", "")
	s.PatchEnvironment("FTP_PROXY", "ftp://user@10.0.0.2")

	proxies := osenv.DetectProxies()

	c.Assert(proxies, gc.DeepEquals, osenv.ProxySettings{
		Http:  "http://user@10.0.0.2",
		Https: "https://user@10.0.0.2",
		Ftp:   "ftp://user@10.0.0.2",
	})
}

func (s *proxySuite) TestDetectPrimaryPreference(c *gc.C) {
	// Patch all of the environment variables we check out just in case the
	// user has one set.
	s.PatchEnvironment("http_proxy", "http://user@10.0.0.1")
	s.PatchEnvironment("https_proxy", "https://user@10.0.0.1")
	s.PatchEnvironment("ftp_proxy", "ftp://user@10.0.0.1")
	s.PatchEnvironment("HTTP_PROXY", "http://user@10.0.0.2")
	s.PatchEnvironment("HTTPS_PROXY", "https://user@10.0.0.2")
	s.PatchEnvironment("FTP_PROXY", "ftp://user@10.0.0.2")

	proxies := osenv.DetectProxies()

	c.Assert(proxies, gc.DeepEquals, osenv.ProxySettings{
		Http:  "http://user@10.0.0.1",
		Https: "https://user@10.0.0.1",
		Ftp:   "ftp://user@10.0.0.1",
	})
}