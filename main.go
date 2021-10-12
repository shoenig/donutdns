package main

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"
	"gophers.dev/cmds/donutdns/agent"
	"gophers.dev/cmds/donutdns/plugins/donutdns"
	"gophers.dev/pkgs/extractors/env"
	"gophers.dev/pkgs/loggy"
)

var directives = []string{
	"startup",
	"debug",
	"log",
	"donutdns",
	"forward",
	"shutdown",
}

func getCC(log loggy.Logger) *agent.CoreConfig {
	cc := agent.ConfigFromEnv(env.OS)
	agent.ApplyDefaults(cc)
	cc.Log(log)
	return cc
}

func init() {
	log := loggy.New("init")
	cc := getCC(log)

	dnsserver.Port = strconv.Itoa(cc.Port)
	dnsserver.Directives = directives
	caddy.SetDefaultCaddyfileLoader(donutdns.PluginName, caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
		return caddy.CaddyfileInput{
			Filepath:       donutdns.PluginName,
			Contents:       []byte(cc.Generate()),
			ServerTypeName: "dns",
		}, nil
	}))
}

func main() {
	coremain.Run()
}
