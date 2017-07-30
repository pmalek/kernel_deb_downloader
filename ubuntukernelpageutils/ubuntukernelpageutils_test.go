package ubuntukernelpageutils

import "testing"

import "strings"
import "reflect"

func equalStringSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

type removeDuplicatesTestData struct {
	str      []string
	expected []string
}

func Test_RemoveDuplicates(t *testing.T) {
	var removeDuplicatesTests = []removeDuplicatesTestData{
		{[]string{"aa", "bb", "cc", "aa", "aa", "aa", "aa", "aa"}, []string{"aa", "bb", "cc"}},
		{[]string{"aa"}, []string{"aa"}},
		{[]string{"aa", "bb", "cc", "aa"}, []string{"aa", "bb", "cc"}},
		{[]string{"aa", "bb", "cc", "aa", "aa", "aa", "123"}, []string{"aa", "bb", "cc", "123"}},
		{[]string{}, []string{}},
	}

	for _, tt := range removeDuplicatesTests {
		actual := removeDuplicates(tt.str)
		if !equalStringSlices(actual, tt.expected) {
			t.Errorf("removeDuplicates(%q): Expected: %q, actual %q", tt.str, tt.expected, actual)
		}
	}
}

type parsePagesTestData struct {
	str      string
	expected map[string]string
}

func Test_parseKernelPage(t *testing.T) {
	var parsePagesTests = []parsePagesTestData{
		{`<tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12.1/">v4.12.1/</a></td><td align="right">2017-07-12 17:20  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041201": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.1/",
			},
		},
		{`<tr><td valign="top"></td><td><a href="v4.12.1/">v4.12.1/</a></td><td align="right">2017-07-12 17:20  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041201": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.1/",
			},
		},
		{`<tr><td><a href="v4.12.1/">v4.12.1/</a></td><td align="right">2017-07-12 17:20  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041201": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.1/",
			},
		},
		{`<tr><td><a href="v4.12.1/">v4.12.1/</a></td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041201": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.1/",
			},
		},
		{`<tr><td><a href="v4.12.1/">v4.12.1/</a></td></tr>`,
			map[string]string{
				"041201": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.1/",
			},
		},
		{`<tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12.4/">v4.12.4/</a></td><td align="right">2017-07-28 01:00  </td><td align="right">  - </td><td>&nbsp;</td></tr><tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.11.10/">v4.11.10/</a></td><td align="right">2017-07-12 16:20  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041110": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.10/",
				"041204": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			},
		},
	}

	for _, tt := range parsePagesTests {
		actual := parseKernelPage(strings.NewReader(tt.str))
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("parseKernelPage(%q)\nExpected: %q,\nactual %q", tt.str, tt.expected, actual)
		}
	}
}

func Test_parseKernelPage_RCsAreNotReturned(t *testing.T) {
	var parsePagesTests = []parsePagesTestData{
		{`<tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12.4/">v4.12.4/</a></td><td align="right">2017-07-28 01:00  </td><td align="right">  - </td><td>&nbsp;</td></tr><tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.11.10/">v4.11.10/</a></td><td align="right">2017-07-12 16:20  </td><td align="right">  - </td><td>&nbsp;</td></tr><tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12-rc3/">v4.12-rc3/</a></td><td align="right">2017-05-29 02:50  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041110": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.10/",
				"041204": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			},
		},
		{`<tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.11.10/">v4.11.10/</a></td><td align="right">2017-07-12 16:20  </td><td align="right">  - </td><td>&nbsp;</td></tr><tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12-rc3/">v4.12-rc3/</a></td><td align="right">2017-05-29 02:50  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{
				"041110": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.11.10/",
			},
		},
		{`<tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="v4.12-rc3/">v4.12-rc3/</a></td><td align="right">2017-05-29 02:50  </td><td align="right">  - </td><td>&nbsp;</td></tr>`,
			map[string]string{},
		},
	}

	for _, tt := range parsePagesTests {
		actual := parseKernelPage(strings.NewReader(tt.str))
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("parseKernelPage(%q)\nExpected: %q,\nactual %q", tt.str, tt.expected, actual)
		}
	}
}
