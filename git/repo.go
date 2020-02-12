package git

import (
	"errors"
	"fmt"
	"github.com/battlemidget/mf/common"

	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

type GitRepo struct {
	session *sh.Session
	tmpDir  string
}

// Sets basic directory level git config settings
func (c *GitRepo) SetGitConfig() {
	c.session.Command("git", "config", "user.email", "cdkbot@gmail.com").Run()
	c.session.Command("git", "config", "user.name", "cdkbot").Run()
	c.session.Command("git", "config", "--global", "push.default", "simple").Run()
}

// clone repo
func CloneRepo(cloneUrl string, tmpDir string) error {
	output, err := sh.Command("git", "clone", "-q", cloneUrl, tmpDir).CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal("Could not clone directory.")
		return errors.New(fmt.Sprintf("Unable to clone repo, %v", err))
	}
	return nil
}

func (c *GitRepo) Checkout(branch string) {
	c.session.Command("git", "checkout", branch, "-q").Run()
}

func (c *GitRepo) Merge(remoteBranch string) {
	c.session.Command("git", "merge", remoteBranch, "-q").Run()
}

func (c *GitRepo) Push(refName string) error {
	_, err := c.session.Command("git", "push", refName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to push to downstream, %v", err))
	}
	return nil
}

func (c *GitRepo) AddRemote(refName string, upstream string) error {
	output, err := c.session.Command("git", "remote", "add", "-f", refName, strings.TrimRight(upstream, "/")).CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal("Could not add remote.")
		return errors.New(fmt.Sprintf("Unable to add remote repo, %v", err))
	}
	return nil

}

// Syncs the upstream repos to our namespace in github
func SyncRepoNamespace(cloneUrl string, repo *common.Repo, dryRun bool) error {
	var c GitRepo
	tmpDir, err := ioutil.TempDir("", "reposync")
	c.tmpDir = tmpDir
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Failed to create tempdir")
		os.Exit(1)
	}
	defer os.RemoveAll(c.tmpDir)

	log.WithFields(log.Fields{"upstream": repo.Upstream, "charm/layer": repo.Name, "dir": c.tmpDir}).Info("Processing repo")

	CloneRepo(cloneUrl, c.tmpDir)

	session := sh.NewSession()
	session.SetDir(c.tmpDir)
	c.session = session

	c.SetGitConfig()
	c.AddRemote("upstream", repo.Upstream)
	c.Checkout("master")
	c.Merge("upstream/master")

	if !dryRun {
		err := c.Push("origin")
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to push: %v", err))
		}
	}
	return nil
}
