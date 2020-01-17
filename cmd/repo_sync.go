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
	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var spec []string
var dryrun bool

type RepoSpec struct {
	Repos []struct {
		Name            string   `yaml:"name"`
		Downstream      string   `yaml:"downstream"`
		Upstream        string   `yaml:"upstream"`
		Tags            []string `yaml:"tags,omitempty"`
		NeedsStable     string   `yaml:"needs_stable,omitempty"`
		NeedsTagging    string   `yaml:"needs_tagging,omitempty"`
		ResourceBuildSh string   `yaml:"resource_build_sh,omitempty"`
		Namespace       string   `yaml:"namespace,omitempty"`
	} `yaml:"repos"`
}

func (c *RepoSpec) parse(specFile string) *RepoSpec {
	filename, _ := filepath.Abs(specFile)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read %s %v", filename, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Failed to unmarshal: %v", err)
		os.Exit(1)
	}
	return c
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync upstream repos",
	Long:  `Syncs upstream git repos with `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Syncing upstream <-> downstream repositories")

		var c RepoSpec
		ghUser := url.QueryEscape(os.Getenv("CDKBOT_GH_USR"))
		ghPass := url.QueryEscape(os.Getenv("CDKBOT_GH_PSW"))
		for _, v := range spec {
			output := c.parse(v)
			for _, r := range output.Repos {
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
					continue
				}
				uPathStrip := strings.TrimRight(uPath.EscapedPath(), ".git")
				uPathStrip = strings.TrimLeft(uPathStrip, "/")
				if r.Downstream == uPathStrip {
					log.WithFields(log.Fields{"downstream": r.Downstream, "upstream": uPathStrip}).Warn("Upstream and downstream are the same, skipping this repository.")
					continue
				}
				cloneUrl := fmt.Sprintf("https://%s:%s@github.com/%s", ghUser, ghPass, r.Downstream)
				if !dryrun {
					// Handle clone
					output, err := sh.Command("git", "clone", cloneUrl, tmpDir).CombinedOutput()
					if err != nil {
						log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal("Could not clone directory.")
						os.Exit(1)
					}

					session := sh.NewSession()
					session.SetDir(tmpDir)
					session.Command("git", "config", "user.email", "cdkbot@juju.solutions").Run()
					session.Command("git", "config", "user.name", "cdkbot").Run()
					session.Command("git", "config", "--global", "push.default", "simple").Run()
					session.Command("git", "remote", "add", "upstream", strings.TrimRight(r.Upstream, "/")).Run()
					session.Command("git", "fetch", "upstream", "-q").Run()
					session.Command("git", "checkout", "master", "-q").Run()
					session.Command("git", "merge", "upstream/master", "-q").Run()
					output, err = session.Command("git", "push", "origin").CombinedOutput()
					if err != nil {
						log.WithFields(log.Fields{"error": err, "downstream": r.Downstream, "output": string(output)}).Fatal("Unable to push to downstream remote")
						os.Exit(1)
					}
				}

			}
		}
	},
}

func init() {
	repoCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringSliceVar(&spec, "spec", []string{}, "Path to YAML specs containing the upstream/downstream locations of the git repos")
	syncCmd.Flags().BoolVar(&dryrun, "dry-run", false, "dry-run")
}
