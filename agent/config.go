package agent

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"gophers.dev/pkgs/extractors/env"
	"gophers.dev/pkgs/loggy"
)

//go:embed config.tmpl
var tmpl embed.FS

type Forward struct {
	Addresses  []string
	ServerName string
}

type CoreConfig struct {
	Port       int
	NoDebug    bool
	NoLog      bool
	Allows     []string
	Blocks     []string
	NoDefaults bool
	Forward    Forward
}

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

func ConfigFromEnv(e env.Environment) *CoreConfig {
	var (
		allow     string
		block     string
		upstream1 string
		upstream2 string
	)

	var cc CoreConfig
	if err := env.Parse(e, env.Schema{
		"DONUT_DNS_PORT":          env.Int(&cc.Port, false),
		"DONUT_DNS_NO_DEBUG":      env.Bool(&cc.NoDebug, false),
		"DONUT_DNS_NO_LOG":        env.Bool(&cc.NoLog, false),
		"DONUT_DNS_ALLOW":         env.String(&allow, false),
		"DONUT_DNS_BLOCK":         env.String(&block, false),
		"DONUT_DNS_NO_DEFAULTS":   env.Bool(&cc.NoDefaults, false),
		"DONUT_DNS_UPSTREAM_1":    env.String(&upstream1, false),
		"DONUT_DNS_UPSTREAM_2":    env.String(&upstream2, false),
		"DONUT_DNS_UPSTREAM_NAME": env.String(&cc.Forward.ServerName, false),
	}); err != nil {
		panic(err)
	}

	var upstreams []string
	if upstream1 != "" {
		upstreams = append(upstreams, upstream1)
	}
	if upstream2 != "" {
		upstreams = append(upstreams, upstream2)
	}
	cc.Forward.Addresses = upstreams

	cc.Allows = split(allow)
	cc.Blocks = split(block)

	return &cc
}

func (cc *CoreConfig) Log(log loggy.Logger) {
	log.Tracef("DONUT_DNS_PORT: %d", cc.Port)
	log.Tracef("DONUT_DNS_DEBUG: %t", cc.NoDebug)
	log.Tracef("DONUT_DNS_NO_LOG: %t", cc.NoLog)
	log.Tracef("DONUT_DNS_NO_ALLOW: %v", cc.Allows)
	log.Tracef("DONUT_DNS_BLOCK: %v", cc.Blocks)
	log.Tracef("DONUT_DNS_NO_DEFAULTS: %t", cc.NoDefaults)
	log.Tracef("DONUT_DNS_UPSTREAM_1: %s", cc.Forward.Addresses[0])
	if len(cc.Forward.Addresses) == 2 {
		log.Tracef("DONUT_DNS_UPSTREAM_2: %s", cc.Forward.Addresses[1])
	}
	log.Tracef("DONUT_DNS_UPSTREAM_NAME: %s", cc.Forward.ServerName)
}

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
