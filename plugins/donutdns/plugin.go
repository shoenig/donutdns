package donutdns

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/shoenig/donutdns/sources/set"
)

const (
	PluginName = "donutdns"
	ServerType = "dns"
)

type DonutDNS struct {
	Next plugin.Handler

	defaultLists bool
	block        *set.Set
	allow        *set.Set
}

func (dd DonutDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	query := state.Name()

	if dd.allow.Has(query) {
		pLog.Infof("query for %s is explicitly allowed", query)
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	if !dd.block.Has(query) {
		pLog.Debugf("query for %s is implicitly allowed", query)
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	var answers []dns.RR

	switch state.QType() {
	case dns.TypeA:
		answers = dd.a(query)
	case dns.TypeAAAA:
		answers = dd.aaaa(query)
	default:
		pLog.Debugf("query: %s type: %d not an A or AAAA record, fallthrough", query, state.QType())
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	qType := dns.Type(state.QType()).String()
	pLog.Infof("BLOCK query (%s) for %s", qType, query)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers
	if err := w.WriteMsg(m); err != nil {
		pLog.Errorf("failed to write msg: %v", err)
		return dns.RcodeServerFailure, err
	}

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
	r := new(dns.AAAA)
	r.Hdr = dns.RR_Header{
		Name:   zone,
		Rrtype: dns.TypeAAAA,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	r.AAAA = sinkAAAA
	return []dns.RR{r}
}
