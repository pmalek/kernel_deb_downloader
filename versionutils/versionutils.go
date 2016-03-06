package versionutils

import (
	"regexp"
	"strings"
)

var reg_dig_in_ver = regexp.MustCompile(`\d+\.?`)
var reg_rc = regexp.MustCompile(`.*rc\d+.*`)
var reg_version = regexp.MustCompile(`v\d+\.\d+\.\d+.*`)

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func leftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func UnifiedVersion(s string, padding int) string {
	var ret string

	versions := reg_dig_in_ver.FindAllString(s, 3)
	for i, v := range versions {
		versions[i] = stripchars(v, ".")
		ret += leftPad2Len(versions[i], "0", padding)
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
