package snap

import (
	// "errors"
	"fmt"
	// "github.com/battlemidget/mf/common"

	"github.com/codeskyblue/go-sh"
	// log "github.com/sirupsen/logrus"
	"strings"
	"regexp"
)

func Revisions(snap string, track string, arch string, excludePre bool) string {
	output, _ := sh.Command("snapcraft", "revisions", snap).Output()
	return string(output)
}


func LatestRev(snap string, arch string, track string, revisions string) int {
	revLines := strings.Split(revisions, "\n")
	revRe := regexp.MustCompile("[ \t+]{2,}")
	var revParse [][]string
	for _, v := range revLines[1:] {
		revParse = append(revParse, revRe.Split(v, -1))
		fmt.Printf("%q\n", revRe.Split(v, -1))
	}

	return 0
}
