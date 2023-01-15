package agent

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/donutdns/output"
	"github.com/shoenig/extractors/env"
)

//go:embed config.tmpl
var tmpl embed.FS

// Forward contains DNS recursion configuration.
type Forward struct {
	Addresses  []string
	ServerName string
}

// CoreConfig contains donutdns configuration.
// It is used to generate CoreDNS (Caddy) style configuration blocks.
type CoreConfig struct {
	Port       int
	NoDebug    bool
	NoLog      bool
	Allows     []string
	AllowFile  string
	AllowDir   string
	Blocks     []string
	BlockFile  string
	BlockDir   string
	Suffix     []string
	SuffixFile string
	SuffixDir  string
	NoDefaults bool
	Forward    *Forward
}

// Generate a CoreDNS (Caddy) style configuration block as a string.
func (cc *CoreConfig) Generate() string {
	t, pErr := template.ParseFS(tmpl, "*.tmpl")
	if pErr != nil {
		panic(pErr)
	}

	var b bytes.Buffer
	if eErr := t.Execute(&b, cc); eErr != nil {
		panic(eErr)
	}
	return b.String()
}

// ConfigFromEnv parses environment variables from e and creates a CoreConfig.
func ConfigFromEnv(e env.Environment) *CoreConfig {
	var (
		allow     string
		block     string
		suffix    string
		upstream1 string
		upstream2 string
	)

	cc := new(CoreConfig)
	cc.Forward = new(Forward)
	if err := env.Parse(e, env.Schema{
		"DONUT_DNS_PORT":          env.Int(&cc.Port, false),
		"DONUT_DNS_NO_DEBUG":      env.Bool(&cc.NoDebug, false),
		"DONUT_DNS_NO_LOG":        env.Bool(&cc.NoLog, false),
		"DONUT_DNS_ALLOW":         env.String(&allow, false),
		"DONUT_DNS_ALLOW_FILE":    env.String(&cc.AllowFile, false),
		"DONUT_DNS_ALLOW_DIR":     env.String(&cc.AllowDir, false),
		"DONUT_DNS_BLOCK":         env.String(&block, false),
		"DONUT_DNS_BLOCK_FILE":    env.String(&cc.BlockFile, false),
		"DONUT_DNS_BLOCK_DIR":     env.String(&cc.BlockDir, false),
		"DONUT_DNS_SUFFIX":        env.String(&suffix, false),
		"DONUT_DNS_SUFFIX_FILE":   env.String(&cc.SuffixFile, false),
		"DONUT_DNS_SUFFIX_DIR":    env.String(&cc.SuffixDir, false),
		"DONUT_DNS_NO_DEFAULTS":   env.Bool(&cc.NoDefaults, false),
		"DONUT_DNS_UPSTREAM_1":    env.String(&upstream1, false),
		"DONUT_DNS_UPSTREAM_2":    env.String(&upstream2, false),
		"DONUT_DNS_UPSTREAM_NAME": env.String(&cc.Forward.ServerName, false),
	}); err != nil {
		panic(err)
	}

	if upstream1 != "" {
		cc.Forward.Addresses = append(cc.Forward.Addresses, upstream1)
	}
	if upstream2 != "" {
		cc.Forward.Addresses = append(cc.Forward.Addresses, upstream2)
	}

	cc.Allows = split(allow)
	cc.Blocks = split(block)
	cc.Suffix = split(suffix)

	return cc
}

// Log cc to plog.
func (cc *CoreConfig) Log(logger output.Info) {
	log.Infof("DONUT_DNS_PORT: %d", cc.Port)
	log.Infof("DONUT_DNS_NO_DEBUG: %t", cc.NoDebug)
	log.Infof("DONUT_DNS_NO_LOG: %t", cc.NoLog)
	log.Infof("DONUT_DNS_ALLOW: %v", cc.Allows)
	log.Infof("DONUT_DNS_ALLOW_FILE: %s", cc.AllowFile)
	log.Infof("DONUT_DNS_ALLOW_DIR: %s", cc.AllowDir)
	log.Infof("DONUT_DNS_BLOCK: %v", cc.Blocks)
	log.Infof("DONUT_DNS_BLOCK_FILE: %s", cc.BlockFile)
	log.Infof("DONUT_DNS_BLOCK_DIR: %s", cc.BlockDir)
	log.Infof("DONUT_DNS_SUFFIX: %v", cc.Suffix)
	log.Infof("DONUT_DNS_SUFFIX_FILE: %s", cc.SuffixFile)
	log.Infof("DONUT_DNS_SUFFIX_DIR: %s", cc.SuffixDir)
	log.Infof("DONUT_DNS_NO_DEFAULTS: %t", cc.NoDefaults)
	log.Infof("DONUT_DNS_UPSTREAM_1: %s", cc.Forward.Addresses[0])
	if len(cc.Forward.Addresses) == 2 {
		log.Infof("DONUT_DNS_UPSTREAM_2: %s", cc.Forward.Addresses[1])
	}
	log.Infof("DONUT_DNS_UPSTREAM_NAME: %s", cc.Forward.ServerName)
}

// ApplyDefaults sets reasonable default config values on a CoreConfig if no value is set.
//
// Port defaults to 5301.
// Forward.Addresses defaults to [1.1.1.1, 1.0.0.1] (cloudflare dns servers)
// Forward.ServerName defaults to cloudflare-dns.com (cloudflare dns servers)
func ApplyDefaults(cc *CoreConfig) {
	if cc.Port == 0 {
		cc.Port = 5301
	}
	if len(cc.Forward.Addresses) == 0 {
		cc.Forward.Addresses = []string{"1.1.1.1", "1.0.0.1"}
		cc.Forward.ServerName = "cloudflare-dns.com"
	}
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
