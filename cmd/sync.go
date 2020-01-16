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
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var spec []string

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
		jww.ERROR.Printf("Failed to read %s %v", filename, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		jww.ERROR.Printf("Failed to unmarshal: %v", err)
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
		fmt.Println("Syncing upstream <-> downstream repositories")

		var c RepoSpec
		ghUser := url.QueryEscape(os.Getenv("CDKBOT_GH_USR"))
		ghPass := url.QueryEscape(os.Getenv("CDKBOT_GH_PSW"))
		for _, v := range spec {
			output := c.parse(v)
			for _, r := range output.Repos {
				tmpDir, err := ioutil.TempDir("", "reposync")
				if err != nil {
					jww.ERROR.Printf("Failed to create tempdir: %v", err)
					os.Exit(1)
				}
				defer os.RemoveAll(tmpDir)
				fmt.Println("- storing repos in ", tmpDir)

				fmt.Printf("Processing: %s\n", r.Name)
				uPath, err := url.Parse(r.Upstream)
				if err != nil {
					jww.ERROR.Printf("[Skipping] Problem reading url: %v", err)
					continue
				}
				uPathStrip := strings.TrimRight(uPath.EscapedPath(), ".git")
				uPathStrip = strings.TrimLeft(uPathStrip, "/")
				if r.Downstream == uPathStrip {
					fmt.Printf("[Skipping] %s == %s\n", r.Downstream, uPathStrip)
					continue
				}
				cloneUrl := fmt.Sprintf("https://%s:%s@github.com/%s", ghUser, ghPass, r.Downstream)
				sh.Command("git", "clone", cloneUrl, tmpDir).Run()

				// Handle clone
				session := sh.NewSession()
				session.SetDir(tmpDir)
				session.Command("git", "config", "user.email", "cdkbot@juju.solutions").Run()
				session.Command("git", "config", "user.name", "cdkbot").Run()
				session.Command("git", "config", "--global", "push.default", "simple").Run()
				session.Command("git", "remote", "add", "upstream", r.Upstream).Run()
				session.Command("git", "fetch", "upstream").Run()
				session.Command("git", "checkout", "master").Run()
				session.Command("git", "merge", "upstream/master").Run()
				session.Command("git", "push", "origin").Run()
			}
		}
	},
}

func init() {
	repoCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringSliceVar(&spec, "spec", []string{}, "Path to YAML specs containing the upstream/downstream locations of the git repos")
}
