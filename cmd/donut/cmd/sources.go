package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed sources.json
var defaults []byte

func sources() (map[string][]string, error) {
	if len(conf.Sources) != 0 {
		return conf.Sources, nil
	}
	var sourceJson []byte
	var err error
	if conf.SourcesFile == "" {
		sourceJson = defaults
	} else {
		sourceJson, err = os.ReadFile(conf.SourcesFile)
		if err != nil {
			return nil, fmt.Errorf("while reading sources file: %w", err)
		}
	}
	s := map[string][]string{}
	err = json.Unmarshal(sourceJson, &s)
	if err != nil {
		return nil, fmt.Errorf("while parsing sources file: %w", err)
	}
	conf.Sources = s
	return conf.Sources, nil
}
