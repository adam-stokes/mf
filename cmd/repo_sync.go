package cmd

import (
	"github.com/battlemidget/mf/common"
	"github.com/battlemidget/mf/git"
	"github.com/korovkin/limiter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"runtime"
)

var spec []string
var dryrun bool

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Repo Sync - syncs upstream repos",
	Long: `Syncs upstream git repos with charmed-kubernetes namespace.

EXAMPLE

    $ mf repo sync --spec path/to/charm-matrix-spec.yml

REQUIREMENTS

- Environment

The following environment variables are required:

CDKBOT_GH_USR - Github username
CDKBOT_GH_PSW - Github password`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Syncing upstream <-> downstream repositories")

		var c common.RepoSpec
		for _, v := range spec {
			output := c.Parse(v)
			limit := limiter.NewConcurrencyLimiter(runtime.GOMAXPROCS(0))
			for _, r := range output.Repos {
				repo := r // NOTE: Need for some reason as `r` was not being properly passed to limiter
				limit.Execute(func() {
					err := git.SyncRepoNamespace(&repo, dryrun)
					if err != nil {
						log.WithFields(log.Fields{"error": err}).Error("Failed to clone repo")
					}
				})
			}
			limit.Wait()
		}
	},
}

func init() {
	repoCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringSliceVar(&spec, "spec", []string{}, "Path to YAML specs containing the upstream/downstream locations of the git repos")
	syncCmd.Flags().BoolVar(&dryrun, "dry-run", false, "dry-run")
}
