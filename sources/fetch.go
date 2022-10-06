package sources

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-set"
	"github.com/shoenig/donutdns/output"
	"github.com/shoenig/donutdns/sources/extract"
	"github.com/shoenig/ignore"
)

// A Downloader is used to download a set of source lists.
type Downloader interface {
	// Download all sources in Lists.
	Download(*Lists) (*set.Set[string], error)
}

type downloader struct {
	logger output.Logger
}

// NewDownloader creates a new Downloader for downloading source lists.
func NewDownloader(logger output.Logger) Downloader {
	return &downloader{
		logger: logger,
	}
}

func (d *downloader) Download(lists *Lists) (*set.Set[string], error) {
	g := NewGetter(d.logger, extract.New(extract.Generic))
	combo := set.New[string](100)
	for _, source := range lists.All() {
		single, err := g.Get(source)
		if err != nil {
			d.logger.Errorf("failed to fetch source %q, skip: %s", source, err)
			continue
		}
		combo.InsertSet(single)
	}
	return combo, nil
}

// A Getter is used to download a single source list.
type Getter interface {
	// Get source and extract its domains into a Set.
	Get(source string) (*set.Set[string], error)
}

type getter struct {
	client *http.Client
	ex     extract.Extractor
	logger output.Logger
}

// NewGetter creates a new Getter, using Extractor ex to extract domains.
func NewGetter(logger output.Logger, ex extract.Extractor) Getter {
	return &getter{
		client: client(
		// todo: pass in one of the upstreams
		//  currently hard-code cloudflare for bootstrapping the sources
		),
		ex:     ex,
		logger: logger,
	}
}

func (g *getter) Get(source string) (*set.Set[string], error) {
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

	g.logger.Infof("got %d domains from %q", single.Size(), source)

	return single, nil
}
