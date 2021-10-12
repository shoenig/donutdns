package donutdns

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"gophers.dev/cmds/donutdns/sources"
	"gophers.dev/cmds/donutdns/sources/set"
)

var plog = log.NewWithPlugin(PluginName)

func init() {
	plugin.Register(PluginName, setup)
}

func setup(c *caddy.Controller) error {

	dd := DonutDNS{
		defaultLists: true, // keep for logging
		block:        set.New(),
		allow:        set.New(),
	}

	for c.Next() {
		_ = c.RemainingArgs()
		for c.NextBlock() {
			switch c.Val() {
			case "defaults":
				if !c.NextArg() {
					return c.ArgErr()
				}
				dd.defaultLists = c.Val() == "true"
				if dd.defaultLists {
					defaults(dd.block)
				}

			case "block":
				if !c.NextArg() {
					return c.ArgErr()
				}
				dd.block.Add(c.Val())

			case "allow":
				if !c.NextArg() {
					return c.ArgErr()
				}
				dd.allow.Add(c.Val())
			}
		}
	}

	plog.Infof("domains allowed: %d", dd.allow.Len())
	plog.Infof("domains blocked: %d", dd.block.Len())

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dd.Next = next
		return dd
	})

	// Plugin loaded okay.
	return nil
}

func defaults(set *set.Set) {
	getter := sources.NewGetter(plog)
	s, err := getter.Get(sources.Defaults())
	if err != nil {
		panic(err)
	}
	set.Union(s)
}
