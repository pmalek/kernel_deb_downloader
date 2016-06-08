package versionutils

import (
	"regexp"
	"strings"

	"github.com/pmalek/stringutils"
)

var reg_dig_in_ver = regexp.MustCompile(`\d+\.?`)
var reg_rc = regexp.MustCompile(`.*rc\d+.*`)
var reg_version = regexp.MustCompile(`v\d+\.\d+\.\d+.*`)

func UnifiedVersion(s string, padding int) string {
	var ret string

	versions := reg_dig_in_ver.FindAllString(s, 3)
	for i, v := range versions {
		versions[i] = stringutils.Stripchars(v, ".")
		ret += stringutils.LeftPad2Len(versions[i], "0", padding)
	}

	ret += strings.Repeat("0", (3-len(versions))*padding)

	return ret
}

func IsAnRCVersion(v string) bool {
	if reg_version.MatchString(v) == false && reg_rc.MatchString(v) == true {
		return true
	}
	return false
}
