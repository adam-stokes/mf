package git

import (
	"errors"
	"fmt"
	"github.com/battlemidget/mf/common"

	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Sets basic directory level git config settings
func SetGitConfig(session *sh.Session) {
	session.Command("git", "config", "user.email", "cdkbot@juju.solutions").Run()
	session.Command("git", "config", "user.name", "cdkbot").Run()
	session.Command("git", "config", "--global", "push.default", "simple").Run()
}

// clone repo
func CloneRepo(cloneUrl string, tmpDir string) error {
	output, err := sh.Command("git", "clone", cloneUrl, tmpDir).CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal("Could not clone directory.")
		return errors.New(fmt.Sprintf("Unable to clone repo, %v", err))
	}
	return nil
}

func Checkout(session *sh.Session, branch string) {
	session.Command("git", "checkout", branch, "-q").Run()
}

func Merge(session *sh.Session, remoteBranch string) {
	session.Command("git", "merge", remoteBranch, "-q").Run()
}

func Push(session *sh.Session, refName string) error {
	_, err := session.Command("git", "push", refName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to push to downstream, %v", err))
	}
	return nil
}

func AddRemote(session *sh.Session, refName string, upstream string) {
	session.Command("git", "remote", "add", "-f", refName, strings.TrimRight(upstream, "/")).Run()
}

// Syncs the upstream repos to our namespace in github
func SyncRepoNamespace(session *sh.Session, cloneUrl string, tmpDir string, repo *common.Repo) error {
	CloneRepo(cloneUrl, tmpDir)
	SetGitConfig(session)
	AddRemote(session, "upstream", repo.Upstream)
	Checkout(session, "master")
	Merge(session, "upstream/master")

	err := Push(session, "origin")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to push: %v", err))
	}
	return nil
}
