package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Allow() (map[string]int, error) {
	return hostFile(conf.Allow)
}

func Block() (map[string]int, error) {
	return hostFile(conf.Block)
}

func hostFile(filename string) (map[string]int, error) {
	if filename == "" {
		return map[string]int{}, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("while reading hostname list file: %w", err)
	}
	entries := map[string]int{}
	scanner := bufio.NewScanner(f)
	lineNo := 0
	regexRe := regexp.MustCompile(Generic)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		host := regexRe.FindString(line)
		if host != "" {
			entries[strings.ToLower(host)] = lineNo
		}
	}
	return entries, nil
}
