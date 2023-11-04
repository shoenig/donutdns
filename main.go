// Command donutdns implements a network level ad-blocking DNS server.
package main

import (
	"context"
	"flag"
	"os"
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	"github.com/coredns/coredns/plugin"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/google/subcommands"
	"github.com/shoenig/donutdns/agent"
	"github.com/shoenig/donutdns/plugins/donutdns"
	"github.com/shoenig/donutdns/subcmds"
	"github.com/shoenig/extractors/env"
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

// pluginLogger is the plugin logger associated with donutdns.
var pluginLogger = log.NewWithPlugin(donutdns.PluginName)

// getCC generates a CoreDNS CoreConfig file using environment variables associated
// with donutdns configuration.
func getCC() *agent.CoreConfig {
	cc := agent.ConfigFromEnv(env.OS)
	agent.ApplyDefaults(cc)
	cc.Log(pluginLogger)
	return cc
}

func setupCC() {
	// get core config from environment
	cc := getCC()

	// set plugin core config
	dnsserver.Port = strconv.Itoa(cc.Port)
	dnsserver.Directives = directives
	caddy.SetDefaultCaddyfileLoader(
		donutdns.PluginName,
		caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
			return caddy.CaddyfileInput{
				Filepath:       donutdns.PluginName,
				Contents:       []byte(cc.Generate()),
				ServerTypeName: donutdns.ServerType,
			}, nil
		}))
}

func main() {
	if len(os.Args) == 1 {
		// launch CoreDNS; plugin configuration must be initialized first
		setupCC()
		plugin.Register(donutdns.PluginName, donutdns.Setup)
		coremain.Run()
		return
	}

	subcommands.Register(subcmds.NewCheckCmd(), "donutdns")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
