package sources

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// client creates an http.Client with an explicit DNS server. This is necessary
// sometimes, like when deploying donutdns on Nomad or Kubernetes where IP tables
// will have 0.0.0.0:53 redirected to this agent, which in turn needs to bootstrap
// itself with the origin sources before it can serve requests.
//
// It might be possible to eliminate this chicken-egg problem using the healthcheck
// plugin or something similar, where we just fallthrough past the donutdns plugin
// until it's actually ready. A project for a rainy day.
//
// Totally ripped from https://koraygocmen.medium.com/custom-dns-resolver-for-the-default-http-client-in-go-a1420db38a5d
func client() *http.Client {
	var (
		dnsResolverIP        = "1.1.1.1:53" // Cloudflare DNS resolver.
		dnsResolverProto     = "udp"        // Protocol to use for the DNS resolver
		dnsResolverTimeoutMs = 5000         // Timeout (ms) for the DNS resolver (optional)
	)

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}).DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	tr := cleanhttp.DefaultTransport()
	tr.DialContext = dialContext

	c := cleanhttp.DefaultClient()
	c.Transport = tr
	return c
}
