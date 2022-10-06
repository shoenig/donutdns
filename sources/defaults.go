package sources

import (
	_ "embed"
	"encoding/json"
)

//go:embed statics/sources.json
var sources []byte

// Lists represents the embedded JSON file containing default source lists.
type Lists struct {
	Suspicious  []string `json:"suspicious"`
	Advertising []string `json:"advertising"`
	Tracking    []string `json:"tracking"`
	Malicious   []string `json:"malicious"`
	Miners      []string `json:"miners"`
}

// Len returns the combined number of default source lists.
func (d *Lists) Len() int {
	return len(d.Suspicious) + len(d.Advertising) + len(d.Tracking) + len(d.Malicious) + len(d.Miners)
}

// All returns a combined list of all default source lists.
func (d *Lists) All() []string {
	all := make([]string, 0, d.Len())
	all = append(all, d.Suspicious...)
	all = append(all, d.Advertising...)
	all = append(all, d.Tracking...)
	all = append(all, d.Malicious...)
	all = append(all, d.Miners...)
	return all
}

// Defaults returns the default set of source lists.
//
// The default set of source lists are embedded as statics/sources.json which
// we then simply unmarshal at runtime.
func Defaults() *Lists {
	lists := new(Lists)
	if err := json.Unmarshal(sources, lists); err != nil {
		panic(err) // defaults are embedded
	}
	return lists
}
