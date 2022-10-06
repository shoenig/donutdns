package donutdns

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/donutdns/agent"
	"github.com/shoenig/donutdns/sources"
)

var pluginLogger = log.NewWithPlugin(PluginName)

// Setup will parse plugin config and register the donutdns plugin
// with the CoreDNS core server.
//
// todo: test with TestController
func Setup(c *caddy.Controller) error {

	// reconstruct the parts of CoreConfig for initializing the allow/block lists
	cc := new(agent.CoreConfig)

	for c.Next() {
		_ = c.RemainingArgs()
		for c.NextBlock() {
			switch c.Val() {
			case "defaults":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.NoDefaults = c.Val() == "false"

			case "allow_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowFile = c.Val()

			case "block_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.BlockFile = c.Val()

			case "suffix_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.SuffixFile = c.Val()

			case "allow":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Allows = append(cc.Allows, c.Val())

			case "block":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Blocks = append(cc.Blocks, c.Val())

			case "suffix":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Suffix = append(cc.Suffix, c.Val())
			}
		}
	}

	sets := sources.New(pluginLogger, cc)
	allow, block, suffix := sets.Size()
	pluginLogger.Infof("domains on explicit allow-list: %d", allow)
	pluginLogger.Infof("domains on explicit block-list: %d", block)
	pluginLogger.Infof("domains on suffixes block-list: %d", suffix)

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dd := DonutDNS{sets: sets}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dd.Next = next
		return dd
	})

	// Plugin loaded okay.
	return nil
}
