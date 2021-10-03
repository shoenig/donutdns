package donutdns

import (
	"gophers.dev/cmds/donutdns/sources/set"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	plugin.Register(PluginName, setup)
}

func setup(c *caddy.Controller) error {

	dd := DonutDNS{
		defaultLists: true,
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

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dd.Next = next
		return dd
	})

	// Plugin loaded okay.
	return nil
}
