package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten",
	Short: "Combine block and allow lists into a single list",
	RunE:  flattenCommand,
}

func init() {
	rootCmd.AddCommand(flattenCmd)
}

func flattenCommand(*cobra.Command, []string) error {
	src, err := sources()
	if err != nil {
		return err
	}
	lists, err := fetchLists(false, false)
	if err != nil {
		return err
	}
	regexRe := regexp.MustCompile(Generic)
	allow, err := Allow()
	if err != nil {
		return err
	}
	block, err := Block()
	if err != nil {
		return err
	}
	fmt.Printf("# Local block list\n")
	for b := range block {
		_, ok := allow[b]
		if !ok {
			fmt.Printf("%s\n", b)
		}
	}
	fmt.Printf("\n\n")
	for cat, urls := range src {
		for _, url := range urls {
			fmt.Printf("\n# %s: %s\n\n", cat, url)
			list := lists[url]
			scanner := bufio.NewScanner(bytes.NewReader(list))
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" || line[0] == '#' {
					continue
				}
				pattern := regexRe.FindString(line)
				host := strings.ToLower(pattern)
				_, ok := allow[host]
				if !ok {
					fmt.Printf("%s\n", host)
				}
			}
		}
	}
	return nil
}
