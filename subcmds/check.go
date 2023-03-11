package subcmds

import (
	"context"
	"flag"
	"strings"

	"github.com/google/subcommands"
	"github.com/shoenig/donutdns/agent"
	"github.com/shoenig/donutdns/output"
	"github.com/shoenig/donutdns/sources"
	"github.com/shoenig/extractors/env"
)

const (
	checkCmdName = "check"
)

type CheckCmd struct {
	quiet    bool
	defaults bool
}

func NewCheckCmd() subcommands.Command {
	return new(CheckCmd)
}
func (cc *CheckCmd) Name() string {
	return checkCmdName
}

func (cc *CheckCmd) Synopsis() string {
	return "Check whether a domain will be blocked."
}

func (cc *CheckCmd) Usage() string {
	return strings.TrimPrefix(`
check [-quiet] [-defaults] <domain>
Check whether domain will be blocked.
`, "\n")
}

func (cc *CheckCmd) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cc.quiet, "quiet", false, "silence verbose debug output")
	fs.BoolVar(&cc.defaults, "defaults", false, "also check against default block lists")
}

func (cc *CheckCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	logger := new(output.CLI)

	args := f.Args()
	if len(args) == 0 {
		logger.Errorf("must specify domain to check command")
		return subcommands.ExitUsageError
	}

	if err := cc.execute(logger, args[0]); err != nil {
		logger.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (cc *CheckCmd) execute(output *output.CLI, domain string) error {
	cfg := agent.ConfigFromEnv(env.OS)
	agent.ApplyDefaults(cfg)
	cfg.NoDefaults = !cc.defaults

	if !cc.quiet {
		cfg.Log(output)
	}

	sets := sources.New(output, cfg)
	switch {
	case sets.Allow(domain):
		output.Infof("domain %q on explicit allow list", domain)
	case sets.BlockByMatch(domain):
		output.Infof("domain %q on explicit block list", domain)
	case sets.BlockBySuffix(domain):
		output.Infof("domain %q on suffix block list", domain)
	default:
		output.Infof("domain %q is implicitly allowable", domain)
	}
	return nil
}
