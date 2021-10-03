package fetch

import (
	"fmt"
	"net/http"

	"gophers.dev/cmds/donutdns/sources/set"

	clean "github.com/hashicorp/go-cleanhttp"
	"gophers.dev/cmds/donutdns/sources/extract"
	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/loggy"
)

type Fetcher interface {
	Fetch(source string) (*set.Set, error)
}

type fetcher struct {
	client *http.Client
	ex     extract.Extractor
	log    loggy.Logger
}

func New(ex extract.Extractor) Fetcher {
	return &fetcher{
		client: clean.DefaultClient(),
		ex:     ex,
		log:    loggy.New("fetch"),
	}
}

func (f *fetcher) Fetch(source string) (*set.Set, error) {
	request, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("User-Agent", "donutDNS")

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

	f.log.Infof("got %d domains from %q", single.Len(), source)

	return single, nil
}
