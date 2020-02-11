/*
Copyright © 2020 Adam Stokes <adam.stokes@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"github.com/battlemidget/mf/common"
	"github.com/battlemidget/mf/git"
	"github.com/codeskyblue/go-sh"
	"github.com/korovkin/limiter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
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

		limit := limiter.NewConcurrencyLimiter(runtime.GOMAXPROCS(0))
		var c common.RepoSpec
		ghUser := url.QueryEscape(os.Getenv("CDKBOT_GH_USR"))
		ghPass := url.QueryEscape(os.Getenv("CDKBOT_GH_PSW"))
		for _, v := range spec {
			output := c.Parse(v)
			for _, r := range output.Repos {
				limit.Execute(func() {
					tmpDir, err := ioutil.TempDir("", "reposync")
					if err != nil {
						log.WithFields(log.Fields{"error": err}).Fatal("Failed to create tempdir")
						os.Exit(1)
					}
					defer os.RemoveAll(tmpDir)
					log.WithFields(log.Fields{"upstream": r.Upstream, "charm/layer": r.Name, "dir": tmpDir}).Info("Processing repo")
					uPath, err := url.Parse(r.Upstream)
					if err != nil {
						log.WithFields(log.Fields{"error": err}).Fatal("Skipping, problem reading url")
						return
					}
					uPathStrip := strings.TrimRight(uPath.EscapedPath(), ".git")
					uPathStrip = strings.TrimLeft(uPathStrip, "/")
					if r.Downstream == uPathStrip {
						log.WithFields(log.Fields{"downstream": r.Downstream, "upstream": uPathStrip}).Warn("Upstream and downstream are the same, skipping this repository.")
						return
					}
					cloneUrl := fmt.Sprintf("https://%s:%s@github.com/%s", ghUser, ghPass, r.Downstream)
					session := sh.NewSession()
					session.SetDir(tmpDir)
					err = git.SyncRepoNamespace(session, cloneUrl, tmpDir, &r, dryrun)
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
