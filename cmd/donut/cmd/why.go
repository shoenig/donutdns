package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

const Generic = `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`

var whyCmd = &cobra.Command{
	Use:   "why [hostname]",
	Short: "Find out where a hostname is blocked",
	RunE:  whyCommand,
}

func init() {
	rootCmd.AddCommand(whyCmd)
}

func whyCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
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
	for _, host := range args {
		cleanHost := strings.ToLower(strings.TrimSuffix(host, "."))
		for cat, urls := range src {
			for _, url := range urls {
				list := lists[url]
				scanner := bufio.NewScanner(bytes.NewReader(list))
				lineNo := 0
				for scanner.Scan() {
					line := scanner.Text()
					lineNo++
					if line == "" || line[0] == '#' {
						continue
					}
					pattern := regexRe.FindString(line)
					if strings.ToLower(pattern) == cleanHost {
						fmt.Printf("%s listed as %s on line %d of %s\n", cleanHost, cat, lineNo, url)
					}
				}
			}
		}
		lineNo, ok := block[cleanHost]
		if ok {
			fmt.Printf("%s listed on line %d of local block list\n", cleanHost, lineNo)
		}
		lineNo, ok = allow[cleanHost]
		if ok {
			fmt.Printf("%s will not be blocked, it is on line %d of local allow list\n", cleanHost, lineNo)
		}
	}
	return nil
}
