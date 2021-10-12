package sources

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/coredns/coredns/plugin/pkg/log"
	"gophers.dev/cmds/donutdns/sources/extract"
	"gophers.dev/cmds/donutdns/sources/fetch"
	"gophers.dev/cmds/donutdns/sources/set"
)

//go:embed statics/sources.json
var sources []byte

type Lists struct {
	Suspicious  []string `json:"suspicious"`
	Advertising []string `json:"advertising"`
	Tracking    []string `json:"tracking"`
	Malicious   []string `json:"malicious"`
	Miners      []string `json:"miners"`
}

func (d *Lists) String() string {
	return fmt.Sprintf(
		"<%d %d %d %d %d>",
		len(d.Suspicious), len(d.Advertising), len(d.Tracking), len(d.Malicious), len(d.Miners),
	)
}

func (d *Lists) Len() int {
	return len(d.Suspicious) + len(d.Advertising) + len(d.Tracking) + len(d.Malicious) + len(d.Miners)
}

func (d *Lists) All() []string {
	all := make([]string, 0, d.Len())
	all = append(all, d.Suspicious...)
	all = append(all, d.Advertising...)
	all = append(all, d.Tracking...)
	all = append(all, d.Malicious...)
	all = append(all, d.Miners...)
	return all
}

func Defaults() *Lists {
	defaults := new(Lists)
	if err := json.Unmarshal(sources, defaults); err != nil {
		panic(err) // defaults are embedded
	}
	return defaults
}

type Getter interface {
	Get(*Lists) (*set.Set, error)
}

type getter struct {
	plog log.P
}

func NewGetter(plog log.P) Getter {
	return &getter{
		plog: plog,
	}
}

func (g *getter) Get(lists *Lists) (*set.Set, error) {
	f := fetch.New(g.plog, extract.New())
	combo := set.New()
	for _, source := range lists.All() {
		single, err := f.Fetch(source)
		if err != nil {
			g.plog.Errorf("failed to fetch source %q, skip: %s", source, err)
			continue
		}
		combo.Union(single)
	}
	return combo, nil
}
