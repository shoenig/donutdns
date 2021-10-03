package donutdns

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"gophers.dev/cmds/donutdns/sources/set"
)

const (
	PluginName = "donutdns"
)

var plog = log.NewWithPlugin(PluginName)

type DonutDNS struct {
	Next plugin.Handler

	defaultLists bool
	block        *set.Set
	allow        *set.Set
}

func (dd DonutDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	plog.Debugf("serve dns was called!, use default list: %t", dd.defaultLists)
	plog.Debugf("blocks: %d", dd.block.Len())
	plog.Debugf("allow: %d", dd.allow.Len())

	// 	for _, q := range r.Question {

	// 	}

	return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
}

func (dd DonutDNS) Name() string {
	return PluginName
}

func (dd DonutDNS) Ready() bool {
	return true
}
