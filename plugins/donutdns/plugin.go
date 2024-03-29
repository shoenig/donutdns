package donutdns

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/shoenig/donutdns/sources"
)

const (
	PluginName = "donutdns"
	ServerType = "dns"
)

type DonutDNS struct {
	Next plugin.Handler
	sets *sources.Sets
}

func (dd DonutDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	origQuery := state.Name()
	cleanQuery := strings.Trim(origQuery, ".")

	if dd.sets.Allow(cleanQuery) {
		pluginLogger.Debugf("query for %s is explicitly allowed", cleanQuery)
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	if dd.sets.BlockByMatch(cleanQuery) {
		pluginLogger.Debugf("query for %s is blocked by match", cleanQuery)
		return dd.null(state.QType(), origQuery, ctx, w, r)
	}

	if dd.sets.BlockBySuffix(cleanQuery) {
		pluginLogger.Debugf("query for %s is blocked by suffix", cleanQuery)
		return dd.null(state.QType(), origQuery, ctx, w, r)
	}

	pluginLogger.Debugf("query for %s is implicitly allowed", cleanQuery)
	return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
}

func (dd DonutDNS) null(qType uint16, query string, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var answers []dns.RR
	queryType := dns.Type(qType).String()

	switch qType {
	case dns.TypeA:
		answers = dd.a(query)
	case dns.TypeAAAA:
		answers = dd.aaaa(query)
	case dns.TypeHTTPS:
		answers = dd.https(query)
	default:
		pluginLogger.Debugf("query: %s type: %s not recognized, fallthrough", query, queryType)
		return plugin.NextOrFailure(dd.Name(), dd.Next, ctx, w, r)
	}

	pluginLogger.Infof("BLOCK query (%s) for %s", queryType, query)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers
	if err := w.WriteMsg(m); err != nil {
		pluginLogger.Errorf("failed to write msg: %v", err)
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

func (dd DonutDNS) https(zone string) []dns.RR {
	return dd.a(zone)
}
