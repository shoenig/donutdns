package main

import (
	"fmt"
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

var directives = []string{
	"startup",
	"debug",
	"log",
	"donutdns",
	"forward",
	"shutdown",
}

var plog = log.NewWithPlugin("donutdns")

func getCC() *agent.CoreConfig {
	cc := agent.ConfigFromEnv(env.OS)
	agent.ApplyDefaults(cc)
	cc.Log(plog)
	return cc
}

func init() {
	cc := getCC()

	fmt.Println(cc.Generate())

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
