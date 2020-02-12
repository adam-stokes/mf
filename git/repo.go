package git

import (
	"errors"
	"fmt"
	"github.com/battlemidget/mf/common"

	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// clone repo
func CloneRepo(upstream string, downstream string, tmpDir string) error {
	ghUser := url.QueryEscape(os.Getenv("CDKBOT_GH_USR"))
	ghPass := url.QueryEscape(os.Getenv("CDKBOT_GH_PSW"))
	uPath, err := url.Parse(upstream)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not read URL, %v", err))
	}
	uPathStrip := strings.TrimRight(uPath.EscapedPath(), ".git")
	uPathStrip = strings.TrimLeft(uPathStrip, "/")
	if downstream == uPathStrip {
		return errors.New(fmt.Sprintf("Upstream and downstream are the same."))
	}
	cloneUrl := fmt.Sprintf("https://%s:%s@github.com/%s", ghUser, ghPass, downstream)

	_, err = sh.Command("git", "clone", "-q", cloneUrl, tmpDir).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to clone repo, %v", err))
	}
	return nil
}

type GitRepo struct {
	session  *sh.Session
	tmpDir   string
	repoSpec *common.Repo
}

// Sets basic directory level git config settings
func (c *GitRepo) SetGitConfig() {
	c.session.Command("git", "config", "user.email", "cdkbot@gmail.com").Run()
	c.session.Command("git", "config", "user.name", "cdkbot").Run()
	c.session.Command("git", "config", "--global", "push.default", "simple").Run()
}

func (c *GitRepo) Checkout(branch string) {
	c.session.Command("git", "checkout", branch, "-q").Run()
}

func (c *GitRepo) Merge(remoteBranch string) error {
	output, err := c.session.Command("git", "merge", remoteBranch, "-q").CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal(fmt.Sprintf("Could not merge into %s.", c.tmpDir))
		return errors.New(fmt.Sprintf("Failed to merge, %v", err))
	}
	return nil

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
		log.WithFields(log.Fields{"error": err, "output": string(output), "upstream": c.repoSpec.Upstream, "downstream": c.repoSpec.Downstream}).Fatal("Could not add remote.")
		return errors.New(fmt.Sprintf("Unable to add remote repo, %v", err))
	}
	return nil

}

// Syncs the upstream repos to our namespace in github
func SyncRepoNamespace(repo *common.Repo, dryRun bool) error {
	var c GitRepo
	c.repoSpec = repo
	tmpDir, err := ioutil.TempDir("", "reposync")
	c.tmpDir = tmpDir
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Failed to create tempdir")
		os.Exit(1)
	}

	defer os.RemoveAll(c.tmpDir)

	log.WithFields(log.Fields{"upstream": repo.Upstream, "charm/layer": repo.Name, "dir": c.tmpDir}).Info("Processing repo")

	err = CloneRepo(repo.Upstream, repo.Downstream, c.tmpDir)
	if err != nil {
		log.WithFields(log.Fields{"upstream": repo.Upstream, "charm/layer": repo.Name, "dir": c.tmpDir}).Warn(fmt.Sprintf("Problem cloning repo, %v, skipping.", err))
		return nil
	}

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
