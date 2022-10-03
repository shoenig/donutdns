package donutdns

import (
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/donutdns/sources"
	"github.com/shoenig/donutdns/sources/extract"
	"github.com/shoenig/donutdns/sources/fetch"
	"github.com/shoenig/donutdns/sources/set"
	"github.com/shoenig/ignore"
)

var pLog = log.NewWithPlugin(PluginName)

func init() {
	plugin.Register(PluginName, setup)
}

// todo: test with TestController
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
				if dd.defaultLists {
					defaults(dd.block)
				}

			case "allow_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				if filename := c.Val(); filename != "" {
					custom(c.Val(), dd.allow)
				}

			case "block_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				if filename := c.Val(); filename != "" {
					custom(c.Val(), dd.block)
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

	pLog.Infof("domains on custom allow-list: %d", dd.allow.Len())
	pLog.Infof("domains on custom block-list: %d", dd.block.Len())

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dd.Next = next
		return dd
	})

	// Plugin loaded okay.
	return nil
}

func defaults(set *set.Set) {
	downloader := fetch.NewDownloader(pLog)
	s, err := downloader.Download(sources.Defaults())
	if err != nil {
		panic(err)
	}
	set.Union(s)
}

func custom(filename string, set *set.Set) {
	// for now, everything uses the generic domain extractor
	ex := extract.New(extract.Generic)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer ignore.Close(f)
	s, err := ex.Extract(f)
	if err != nil {
		panic(err)
	}
	set.Union(s)
}
