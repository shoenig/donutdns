package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Config struct {
	SourcesFile   string
	Sources       map[string][]string
	CacheDir      string
	CacheLifetime time.Duration
	Version       string
}

var conf = Config{
	Sources:       map[string][]string{},
	CacheDir:      "",
	CacheLifetime: 3600 * time.Second,
	Version:       "0.1",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "donut",
	Short: "Manage donutdns",
	Long:  `donut is a tool to manage donutdns block lists.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&conf.CacheDir, "cache", "", "cache directory (default is $HOME/.donut)")
	rootCmd.PersistentFlags().StringVar(&conf.SourcesFile, "sources", "", "json sources file (default is built-in defaults)")
}
