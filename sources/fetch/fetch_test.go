package fetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/test/must"
	"gophers.dev/cmds/donutdns/sources"
	"gophers.dev/cmds/donutdns/sources/extract"
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
	g := NewGetter(pLog, ex)
	s, err := g.Get(ts.URL)
	must.NoError(t, err)
	must.EqOp(t, 3, s.Len())
}

func Test_Download(t *testing.T) {
	hit := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, example)
		hit++
	}))
	defer ts.Close()

	lists := &sources.Lists{
		Suspicious:  []string{ts.URL},
		Advertising: []string{ts.URL},
		Tracking:    []string{ts.URL},
		Malicious:   []string{ts.URL},
		Miners:      []string{ts.URL},
	}

	d := NewDownloader(pLog)
	s, err := d.Download(lists)
	must.NoError(t, err)
	must.EqOp(t, 3, s.Len())
	must.EqOp(t, 5, hit)
}
