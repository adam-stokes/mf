package common

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)


type RepoSpec struct {
	Repos []Repo `yaml:"repos"`
}

type Repo struct {
	Name            string   `yaml:"name"`
		Downstream      string   `yaml:"downstream"`
		Upstream        string   `yaml:"upstream"`
		Tags            []string `yaml:"tags,omitempty"`
		NeedsStable     string   `yaml:"needs_stable,omitempty"`
		NeedsTagging    string   `yaml:"needs_tagging,omitempty"`
		ResourceBuildSh string   `yaml:"resource_build_sh,omitempty"`
		Namespace       string   `yaml:"namespace,omitempty"`

}

func (c *RepoSpec) Parse(specFile string) *RepoSpec {
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
