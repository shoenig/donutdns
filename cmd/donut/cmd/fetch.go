package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch upstream blocklists",
	RunE:  fetchCommand,
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

// filenamesFromUrls makes nice filenames to store fetched sources
// in the cache, so they're pleasant-ish to use by hand.
func filenamesFromUrls(sources map[string][]string) map[string]string {
	urlRe := regexp.MustCompile(`^https?://([^/]+)/([^?]+)`)
	ret := map[string]string{}
	for cat, list := range sources {
		seen := map[string]struct{}{}
		var fname string
		for _, url := range list {
			matches := urlRe.FindStringSubmatch(url)
			if matches == nil {
				_, fname = path.Split(url)
			} else {
				hostname := matches[1]
				parts := strings.Split(matches[2], "/")
				switch hostname {
				case "raw.githubusercontent.com", "bitbucket.org", "s3.amazonaws.com":
					hostname = parts[0]
					parts = parts[1:]
				default:
				}
				if len(parts) > 0 {
					fname = hostname + "-" + parts[len(parts)-1]
				} else {
					fname = hostname
				}
			}
			if fname == "" {
				fname = "list"
			}
			_, ok := seen[fname]
			if ok {
				i := 1
				for {
					fn := fmt.Sprintf("%s.%d", fname, i)
					_, ok = seen[fn]
					if !ok {
						fname = fn
						break
					}
					i++
				}
			}
			seen[fname] = struct{}{}
			ret[url] = cat + "-" + fname
		}
	}
	return ret
}

func fetchCommand(*cobra.Command, []string) error {
	_, err := fetchLists(true, true)
	return err
}

func cacheDir() (string, error) {
	cd := conf.CacheDir
	if cd == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cd = filepath.Join(home, ".donutcache")
	}
	err := os.MkdirAll(cd, 0755)
	if err != nil {
		log.Printf("Failed to create '%s': %v\n", cd, err)
		cd, err = os.MkdirTemp(os.TempDir(), "donut-*")
		if err != nil {
			return "", err
		}
		log.Printf("Using '%s' instead\n", cd)
	}
	return cd, nil
}

func fetchLists(force, verbose bool) (map[string][]byte, error) {
	sourceList, err := sources()
	if err != nil {
		return nil, err
	}
	filenames := filenamesFromUrls(sourceList)
	content := map[string][]byte{}
	cache, err := cacheDir()
	if err != nil {
		return nil, err
	}
	for url, file := range filenames {
		filename := filepath.Join(cache, file)
		if !force {
			fi, err := os.Stat(filename)
			if err == nil && fi.ModTime().Add(conf.CacheLifetime).After(time.Now()) {
				cachedContent, err := os.ReadFile(filename)
				if err == nil {
					content[url] = cachedContent
					if verbose {
						log.Printf("%s (cached)", url)
					}
					continue
				}
			}
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("while creating request from %s: %w", url, err)
		}
		req.Header.Set("User-Agent", fmt.Sprintf("com.wordtothewise.donut (%s)", conf.Version))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("`%s (failed, %v)", url, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("%s (failed, %s)", url, resp.Status)
			_ = resp.Body.Close()
			continue
		}
		newContent, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			log.Printf("%s (failed, %v)", url, err)
			continue
		}
		content[filename] = newContent
		err = os.WriteFile(filename, newContent, 0644)
		if err != nil {
			log.Printf("failed to write cache file: %v", err)
			_ = os.Remove(filename)
			continue
		}
		if verbose {
			log.Printf("%s (OK)", url)
		}
	}
	return content, nil
}
