package fetch

import (
	"fmt"
	"net/http"

	"github.com/coredns/coredns/plugin/pkg/log"
	clean "github.com/hashicorp/go-cleanhttp"
	"gophers.dev/cmds/donutdns/sources/extract"
	"gophers.dev/cmds/donutdns/sources/set"
	"gophers.dev/pkgs/ignore"
)

type Fetcher interface {
	Fetch(source string) (*set.Set, error)
}

type fetcher struct {
	client *http.Client
	ex     extract.Extractor
	plog   log.P
}

func New(plog log.P, ex extract.Extractor) Fetcher {
	return &fetcher{
		client: clean.DefaultClient(),
		ex:     ex,
		plog:   plog,
	}
}

func (f *fetcher) Fetch(source string) (*set.Set, error) {
	request, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("User-Agent", "donutdns")

	response, err := f.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer ignore.Drain(response.Body)

	if code := response.StatusCode; code >= 400 {
		return nil, fmt.Errorf("unexpected request response, code: %d", code)
	}

	single, err := f.ex.Extract(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sources: %w", err)
	}

	f.plog.Infof("got %d domains from %q", single.Len(), source)

	return single, nil
}
