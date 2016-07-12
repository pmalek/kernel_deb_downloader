package versionutils

import (
	"regexp"
	"strings"

	"github.com/pmalek/stringutils"
)

var regDigInVer = regexp.MustCompile(`\d+\.?`)
var regRC = regexp.MustCompile(`.*rc\d+.*`)
var regVersion = regexp.MustCompile(`v\d+\.\d+\.\d+.*`)

// UnifiedVersion returns a unified version string where each of
// major, minor and patch parts of version string @s will have a @padding
// number of characters without stripped of characters that are not digits
// e.g. UnifiedVersion("004.004.004", 4) will return "000400040004"
func UnifiedVersion(s string, padding int) string {
	var ret string

	versions := regDigInVer.FindAllString(s, 3)
	for i, v := range versions {
		versions[i] = stringutils.Stripchars(v, ".")
		ret += stringutils.LeftPad2Len(versions[i], "0", padding)
	}

	ret += strings.Repeat("0", (3-len(versions))*padding)

	return ret
}

// IsAnRCVersion return a bool indicating whether @v is an RC version
// e.g. returns true for "v4.6-rc7-wily", return false for "v4.5.0"
func IsAnRCVersion(v string) bool {
	if regVersion.MatchString(v) == false && regRC.MatchString(v) == true {
		return true
	}
	return false
}
