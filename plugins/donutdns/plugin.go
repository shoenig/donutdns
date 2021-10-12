package donutdns

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
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

func allow(name string) bool {
	return name != "facebook.com."
}

func (dd DonutDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	plog.Debugf("serve dns was called!, use default list: %t, blocks: %d, allows: %d", dd.defaultLists, dd.block.Len(), dd.allow.Len())

	// todo: decide on fallthrough

	state := request.Request{W: w, Req: r}
	query := state.Name()

	if allow(query) {
		plog.Debugf("query for %s is not blocked", query)
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	var answers []dns.RR

	plog.Debugf("query: %s, qtype: %d", query, state.QType())

	switch state.QType() {
	case dns.TypeA:
		answers = dd.a(query)
	case dns.TypeAAAA:
		answers = dd.aaaa(query)
	default:
		plog.Debugf("not an A or AAAA record, fallthrough")
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers
	if err := w.WriteMsg(m); err != nil {
		plog.Debugf("failed to write msg: %v", err)
		return dns.RcodeServerFailure, err
	}

	plog.Debugf("wrote response")
	return dns.RcodeSuccess, nil
}

func (dd DonutDNS) Name() string {
	return PluginName
}

func (dd DonutDNS) Ready() bool {
	return true
}

var sinkA = net.IP([]byte{0, 0, 0, 0})
var sinkAAAA = net.IP([]byte{
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
})

func (dd DonutDNS) a(zone string) []dns.RR {
	r := new(dns.A)
	r.Hdr = dns.RR_Header{
		Name:   zone,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	r.A = sinkA
	return []dns.RR{r}
}

func (dd DonutDNS) aaaa(zone string) []dns.RR {
	r := new(dns.A)
	r.Hdr = dns.RR_Header{
		Name:   zone,
		Rrtype: dns.TypeAAAA,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	r.A = sinkAAAA
	return []dns.RR{r}
}
