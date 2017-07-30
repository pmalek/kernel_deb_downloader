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

type parsePackangePageTestData struct {
	str        string
	packageURL string
	expected   []string
}

func Test_parsePackagePage(t *testing.T) {
	var parsePagesTests = []parsePackangePageTestData{
		{str: `<tr><td valign="top"><img src="/icons/back.gif" alt="[PARENTDIR]"></td><td><a href="/~kernel-ppa/mainline/">Parent Directory</a></td><td>&nbsp;</td><td align="right">  - </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/text.gif" alt="[TXT]"></td><td><a href="0001-base-packaging.patch">0001-base-packaging.patch</a></td><td align="right">2017-07-27 23:32  </td><td align="right"> 14M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/text.gif" alt="[TXT]"></td><td><a href="0002-debian-changelog.patch">0002-debian-changelog.patch</a></td><td align="right">2017-07-27 23:32  </td><td align="right"> 40K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/text.gif" alt="[TXT]"></td><td><a href="0003-configs-based-on-Ubuntu-4.12.0-8.9.patch">0003-configs-based-on-Ubuntu-4.12.0-8.9.patch</a></td><td align="right">2017-07-27 23:32  </td><td align="right"> 11K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG">BUILD.LOG</a></td><td align="right">2017-07-28 00:51  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.amd64">BUILD.LOG.amd64</a></td><td align="right">2017-07-28 00:51  </td><td align="right">2.9M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.arm64">BUILD.LOG.arm64</a></td><td align="right">2017-07-28 00:51  </td><td align="right">1.4M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.armhf">BUILD.LOG.armhf</a></td><td align="right">2017-07-28 00:51  </td><td align="right">2.9M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.binary-headers">BUILD.LOG.binary-headers</a></td><td align="right">2017-07-28 00:51  </td><td align="right"> 15K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.i386">BUILD.LOG.i386</a></td><td align="right">2017-07-28 00:51  </td><td align="right">2.9M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.ppc64el">BUILD.LOG.ppc64el</a></td><td align="right">2017-07-28 00:51  </td><td align="right">1.3M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILD.LOG.s390x">BUILD.LOG.s390x</a></td><td align="right">2017-07-28 00:51  </td><td align="right">351K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="BUILT">BUILT</a></td><td align="right">2017-07-28 00:51  </td><td align="right">232 </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="CHANGES">CHANGES</a></td><td align="right">2017-07-27 23:32  </td><td align="right"> 14K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="CHECKSUMS">CHECKSUMS</a></td><td align="right">2017-07-28 00:51  </td><td align="right">5.5K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="CHECKSUMS.gpg">CHECKSUMS.gpg</a></td><td align="right">2017-07-28 01:00  </td><td align="right">473 </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="COMMIT">COMMIT</a></td><td align="right">2017-07-28 00:51  </td><td align="right"> 49 </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/text.gif" alt="[TXT]"></td><td><a href="HEADER.html">HEADER.html</a></td><td align="right">2017-07-28 00:51  </td><td align="right">5.5K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/hand.right.gif" alt="[    ]"></td><td><a href="README">README</a></td><td align="right">2017-07-28 00:51  </td><td align="right">2.7K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="REBUILD">REBUILD</a></td><td align="right">2017-07-28 00:51  </td><td align="right"> 33 </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="SOURCES">SOURCES</a></td><td align="right">2017-07-28 00:51  </td><td align="right">234 </td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb">linux-headers-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:28  </td><td align="right">707K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right">654K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb</a></td><td align="right">2017-07-28 00:38  </td><td align="right">684K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:26  </td><td align="right">712K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:08  </td><td align="right">646K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right">899K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:51  </td><td align="right">350K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51  </td><td align="right">656K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right">645K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb">linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb</a></td><td align="right">2017-07-27 23:33  </td><td align="right"> 10M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb">linux-image-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:28  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb</a></td><td align="right">2017-07-28 00:38  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:26  </td><td align="right"> 48M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:07  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right"> 44M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:50  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51 </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>`,
			packageURL: "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			expected: []string{
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb"},
		},
		{str: `<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb">linux-headers-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:28  </td><td align="right">707K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right">654K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb</a></td><td align="right">2017-07-28 00:38  </td><td align="right">684K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:26  </td><td align="right">712K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:08  </td><td align="right">646K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right">899K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:51  </td><td align="right">350K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51  </td><td align="right">656K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right">645K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb">linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb</a></td><td align="right">2017-07-27 23:33  </td><td align="right"> 10M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb">linux-image-4.12.4-041204-generic-lpae_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:28  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_arm64.deb</a></td><td align="right">2017-07-28 00:38  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_armhf.deb</a></td><td align="right">2017-07-28 00:26  </td><td align="right"> 48M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:07  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right"> 44M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:50  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51 </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>`,
			packageURL: "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			expected: []string{
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb"},
		},
		{str: `<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right">654K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:08  </td><td align="right">646K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right">899K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:51  </td><td align="right">350K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51  </td><td align="right">656K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right">645K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb">linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb</a></td><td align="right">2017-07-27 23:33  </td><td align="right"> 10M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:07  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right"> 44M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:50  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51 </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>`,
			packageURL: "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			expected: []string{
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-headers-4.12.4-041204_4.12.4-041204.201707271932_all.deb",
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb"},
		},
		{str: `<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:08  </td><td align="right">646K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right">899K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:51  </td><td align="right">350K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51  </td><td align="right">656K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right">645K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:49  </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:07  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right"> 44M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:50  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51 </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>`,
			packageURL: "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			expected: []string{
				"http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_amd64.deb"},
		},
		{str: `<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:08  </td><td align="right">646K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right">899K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-headers-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:51  </td><td align="right">350K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51  </td><td align="right">656K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-headers-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right">645K</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:07  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_ppc64el.deb</a></td><td align="right">2017-07-28 00:47  </td><td align="right"> 44M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb">linux-image-4.12.4-041204-generic_4.12.4-041204.201707271932_s390x.deb</a></td><td align="right">2017-07-28 00:50  </td><td align="right"> 12M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_amd64.deb</a></td><td align="right">2017-07-27 23:51 </td><td align="right"> 49M</td><td>&nbsp;</td></tr>
		<tr><td valign="top"><img src="/icons/unknown.gif" alt="[    ]"></td><td><a href="linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb">linux-image-4.12.4-041204-lowlatency_4.12.4-041204.201707271932_i386.deb</a></td><td align="right">2017-07-28 00:10  </td><td align="right"> 47M</td><td>&nbsp;</td></tr>`,
			packageURL: "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.12.4/",
			expected:   []string{},
		},
	}

	for _, tt := range parsePagesTests {
		actual := parsePackagePage(strings.NewReader(tt.str), tt.packageURL)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("parsePackagePage(%q)\nExpected: %q,\nactual %q", tt.str, tt.expected, actual)
		}
	}
}

type getMostActualKernelVersionTestData struct {
	links           map[string]string
	expectedVersion string
	expectedLink    string
}

func Test_getMostActualKernelVersion(t *testing.T) {
	var getMostActualKernelVersionTests = []getMostActualKernelVersionTestData{
		{
			links: map[string]string{
				"040116": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.1.16-wily/",
				"040919": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.9.19/",
				"041015": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.10.15/",
				"040113": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.1.13-wily/",
				"040815": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.8.15/",
			},
			expectedVersion: "041015",
			expectedLink:    "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.10.15/",
		},
		{
			links: map[string]string{
				"040113": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.1.13-wily/",
				"040815": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.8.15/",
			},
			expectedVersion: "040815",
			expectedLink:    "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.8.15/",
		},
		{
			links: map[string]string{
				"040815": "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.8.15/",
			},
			expectedVersion: "040815",
			expectedLink:    "http://kernel.ubuntu.com/~kernel-ppa/mainline/v4.8.15/",
		},
		{
			links:           map[string]string{},
			expectedVersion: "",
			expectedLink:    "",
		},
	}

	for _, tt := range getMostActualKernelVersionTests {
		actualVersion, actualLink := getMostActualKernelVersion(tt.links)
		if actualVersion != tt.expectedVersion || actualLink != tt.expectedLink {
			t.Errorf(" getMostActualKernelVersion(%q)\nExpected: %q, %q,\nactual %q, %q",
				tt.links, tt.expectedVersion, tt.expectedLink, actualVersion, actualLink)
		}
	}
}
