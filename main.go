// Command donutdns implements a network level ad-blocking DNS server.
package main

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"
	"github.com/coredns/coredns/plugin/pkg/log"
	"gophers.dev/cmds/donutdns/agent"
	"gophers.dev/cmds/donutdns/plugins/donutdns"
	"gophers.dev/pkgs/extractors/env"
)

// directives contains the ordered set of plugins to enable in CoreDNS.
var directives = []string{
	"startup",
	"debug",
	"log",
	"donutdns",
	"forward",
	"shutdown",
}

// pLog is the plugin logger associated with donutdns.
var pLog = log.NewWithPlugin(donutdns.PluginName)

// getCC generates a CoreDNS CoreConfig file using environment variables associated with
// donutdns configuration.
func getCC() *agent.CoreConfig {
	cc := agent.ConfigFromEnv(env.OS)
	agent.ApplyDefaults(cc)
	cc.Log(pLog)
	return cc
}

func init() {
	// get core config from environment
	cc := getCC()

	// set plugin core config
	dnsserver.Port = strconv.Itoa(cc.Port)
	dnsserver.Directives = directives
	caddy.SetDefaultCaddyfileLoader(donutdns.PluginName, caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
		return caddy.CaddyfileInput{
			Filepath:       donutdns.PluginName,
			Contents:       []byte(cc.Generate()),
			ServerTypeName: donutdns.ServerType,
		}, nil
	}))
}

func main() {
	// launch CoreDNS; plugin configuration must be in init blocks
	coremain.Run()
}
