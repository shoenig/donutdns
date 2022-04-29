package agent

import (
	"strings"
	"testing"

	"github.com/shoenig/test/must"
	"gophers.dev/pkgs/extractors/env"
)

func TestCoreConfig_Generate(t *testing.T) {
	cc := CoreConfig{
		Port:      1053,
		NoDebug:   true,
		NoLog:     true,
		Allows:    []string{"example.com", "pets.com"},
		AllowFile: "/etc/allow.list",
		Blocks:    []string{"facebook.com", "instagram.com"},
		BlockFile: "/etc/block.list",
		Forward: Forward{
			Addresses:  []string{"1.1.1.1", "1.0.0.1"},
			ServerName: "cloudflare-dns.com",
		},
	}

	result := cc.Generate()
	must.EqOp(t, noWhitespace(`
.:1053 {
  donutdns {
    defaults true
    allow_file /etc/allow.list
    block_file /etc/block.list

    allow example.com
    allow pets.com
    
    block facebook.com
    block instagram.com
    
  }
  forward . 1.1.1.1 1.0.0.1 {
    tls_servername cloudflare-dns.com
  }
}
`), noWhitespace(result))
}

func TestCoreConfig_Generate_less(t *testing.T) {
	cc := CoreConfig{
		Port:       1054,
		NoDebug:    false,
		NoLog:      false,
		Allows:     nil,
		Blocks:     nil,
		NoDefaults: true,
		Forward: Forward{
			Addresses:  []string{"8.8.8.8"},
			ServerName: "google.dns",
		},
	}

	result := cc.Generate()
	must.EqOp(t, noWhitespace(`
.:1054 {
  debug
  log
  donutdns {
    defaults false
  }
  forward . 8.8.8.8 {
    tls_servername google.dns
  }
}
`), noWhitespace(result))
}

func noWhitespace(s string) string {
	a := strings.ReplaceAll(s, " ", "")
	b := strings.ReplaceAll(a, "\n", "")
	return b
}

func TestConfigFromEnv(t *testing.T) {
	mEnv := env.NewEnvironmentMock(t)
	defer mEnv.MinimockFinish()

	mEnv.GetenvMock.When("DONUT_DNS_PORT").Then("1234")
	mEnv.GetenvMock.When("DONUT_DNS_NO_DEBUG").Then("1")
	mEnv.GetenvMock.When("DONUT_DNS_NO_LOG").Then("1")
	mEnv.GetenvMock.When("DONUT_DNS_ALLOW").Then("example.com,pets.com")
	mEnv.GetenvMock.When("DONUT_DNS_ALLOW_FILE").Then("/etc/allow.list")
	mEnv.GetenvMock.When("DONUT_DNS_BLOCK").Then("facebook.com,reddit.com")
	mEnv.GetenvMock.When("DONUT_DNS_BLOCK_FILE").Then("/etc/block.list")
	mEnv.GetenvMock.When("DONUT_DNS_NO_DEFAULTS").Then("")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_1").Then("8.8.8.8")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_2").Then("8.8.4.4")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_NAME").Then("dns.google")

	cc := ConfigFromEnv(mEnv)
	must.Eq(t, &CoreConfig{
		Port:       1234,
		NoDebug:    true,
		NoLog:      true,
		Allows:     []string{"example.com", "pets.com"},
		AllowFile:  "/etc/allow.list",
		Blocks:     []string{"facebook.com", "reddit.com"},
		BlockFile:  "/etc/block.list",
		NoDefaults: false,
		Forward: Forward{
			Addresses:  []string{"8.8.8.8", "8.8.4.4"},
			ServerName: "dns.google",
		},
	}, cc)
}

func TestConfigFromEnv_2(t *testing.T) {
	mEnv := env.NewEnvironmentMock(t)
	defer mEnv.MinimockFinish()

	mEnv.GetenvMock.When("DONUT_DNS_PORT").Then("1234")
	mEnv.GetenvMock.When("DONUT_DNS_NO_DEBUG").Then("0")
	mEnv.GetenvMock.When("DONUT_DNS_NO_LOG").Then("true")
	mEnv.GetenvMock.When("DONUT_DNS_ALLOW").Then("")
	mEnv.GetenvMock.When("DONUT_DNS_ALLOW_FILE").Then("")
	mEnv.GetenvMock.When("DONUT_DNS_BLOCK").Then("facebook.com")
	mEnv.GetenvMock.When("DONUT_DNS_BLOCK_FILE").Then("")
	mEnv.GetenvMock.When("DONUT_DNS_NO_DEFAULTS").Then("true")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_1").Then("8.8.8.8")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_2").Then("")
	mEnv.GetenvMock.When("DONUT_DNS_UPSTREAM_NAME").Then("dns.google")

	cc := ConfigFromEnv(mEnv)
	must.Eq(t, &CoreConfig{
		Port:       1234,
		NoDebug:    false,
		NoLog:      true,
		Allows:     nil,
		Blocks:     []string{"facebook.com"},
		BlockFile:  "",
		NoDefaults: true,
		Forward: Forward{
			Addresses:  []string{"8.8.8.8"},
			ServerName: "dns.google",
		},
	}, cc)
}
