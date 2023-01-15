package sources

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/donutdns/agent"
	"github.com/shoenig/donutdns/sources/extract"
	"github.com/shoenig/test/must"
)

var pLog = log.NewWithPlugin("-test")

const example = `
# [socials]
facebook.com
instagram.com
twitter.com`

func Test_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, example)
	}))
	defer ts.Close()

	ex := extract.New(extract.Generic)
	fwd := new(agent.Forward)

	g := NewGetter(pLog, fwd, ex)
	s, err := g.Get(ts.URL)
	must.NoError(t, err)
	must.EqOp(t, 3, s.Size())
}

func Test_Get_bad_upstream(t *testing.T) {
	ex := extract.New(extract.Generic)
	fwd := &agent.Forward{Addresses: []string{"0.0.0.0"}}

	g := NewGetter(pLog, fwd, ex)
	_, err := g.Get("http://example.com")
	must.ErrorContains(t, err, "dial tcp: lookup example.com")
}

func Test_Download(t *testing.T) {
	hit := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, example)
		hit++
	}))
	defer ts.Close()

	lists := &Lists{
		Suspicious:  []string{ts.URL},
		Advertising: []string{ts.URL},
		Tracking:    []string{ts.URL},
		Malicious:   []string{ts.URL},
		Miners:      []string{ts.URL},
	}

	fwd := new(agent.Forward)
	d := NewDownloader(fwd, pLog)

	s, err := d.Download(lists)
	must.NoError(t, err)
	must.EqOp(t, 3, s.Size())
	must.EqOp(t, 5, hit)
}
