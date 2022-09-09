package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

//go:embed sources.json
var defaults []byte

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "print the embedded list of sources",
	RunE:  defaultsCommand,
}

func init() {
	rootCmd.AddCommand(defaultsCmd)
}

func defaultsCommand(*cobra.Command, []string) error {
	_, err := os.Stdout.Write(defaults)
	return err
}

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
