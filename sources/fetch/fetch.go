package fetch

import (
	"fmt"
	"net/http"

	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/shoenig/donutdns/sources"
	"github.com/shoenig/donutdns/sources/extract"
	"github.com/shoenig/donutdns/sources/set"
	"github.com/shoenig/ignore"
)

// A Downloader is used to download a set of source lists.
type Downloader interface {
	// Download all sources in Lists.
	Download(*sources.Lists) (*set.Set, error)
}

type downloader struct {
	pLog log.P
}

// NewDownloader creates a new Downloader for downloading source lists.
func NewDownloader(pLog log.P) Downloader {
	return &downloader{
		pLog: pLog,
	}
}

func (d *downloader) Download(lists *sources.Lists) (*set.Set, error) {
	g := NewGetter(d.pLog, extract.New(extract.Generic))
	combo := set.New()
	for _, source := range lists.All() {
		single, err := g.Get(source)
		if err != nil {
			d.pLog.Errorf("failed to fetch source %q, skip: %s", source, err)
			continue
		}
		combo.Union(single)
	}
	return combo, nil
}

// A Getter is used to download a single source list.
type Getter interface {
	// Get source and extract its domains into a Set.
	Get(source string) (*set.Set, error)
}

type getter struct {
	client *http.Client
	ex     extract.Extractor
	plog   log.P
}

// NewGetter creates a new Getter, using Extractor ex to extract domains.
func NewGetter(pLog log.P, ex extract.Extractor) Getter {
	return &getter{
		client: client(
		// todo: pass in one of the upstreams
		//  currently hard-code cloudflare for bootstrapping the sources
		),
		ex:   ex,
		plog: pLog,
	}
}

func (g *getter) Get(source string) (*set.Set, error) {
	request, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("User-Agent", "donutdns")

	response, err := g.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer ignore.Drain(response.Body)

	if code := response.StatusCode; code >= 400 {
		return nil, fmt.Errorf("unexpected request response, code: %d", code)
	}

	single, err := g.ex.Extract(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sources: %w", err)
	}

	g.plog.Infof("got %d domains from %q", single.Len(), source)

	return single, nil
}
