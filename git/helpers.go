package git

import (
	"errors"
	"fmt"
	"github.com/battlemidget/mf/common"

	"github.com/codeskyblue/go-sh"
	log "github.com/sirupsen/logrus"
	"strings"
)

func CloneRepo(session *sh.Session, cloneUrl string, tmpDir string, repo *common.Repo) error {
	output, err := sh.Command("git", "clone", cloneUrl, tmpDir).CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "output": string(output)}).Fatal("Could not clone directory.")
		return errors.New(fmt.Sprintf("Unable to clone repo, %v", err))
	}

	session.Command("git", "config", "user.email", "cdkbot@juju.solutions").Run()
	session.Command("git", "config", "user.name", "cdkbot").Run()
	session.Command("git", "config", "--global", "push.default", "simple").Run()
	session.Command("git", "remote", "add", "upstream", strings.TrimRight(repo.Upstream, "/")).Run()
	session.Command("git", "fetch", "upstream", "-q").Run()
	session.Command("git", "checkout", "master", "-q").Run()
	session.Command("git", "merge", "upstream/master", "-q").Run()
	output, err = session.Command("git", "push", "origin").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to push to downstream, %v", err))
	}
	return nil
}
