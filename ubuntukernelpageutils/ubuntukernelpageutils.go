package ubuntukernelpageutils

import (
	"./../downloadutils"
	"./../versionutils"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)

var Kernel_webpage = "http://kernel.ubuntu.com/~kernel-ppa/mainline/"

func parseKernelPage() (links map[string]string) {
	const padding = 2

	resp, _ := http.Get(Kernel_webpage)
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	links = make(map[string]string, 0) // the same as []string{}

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken: // End of the document, we're done
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" && versionutils.IsAnRCVersion(a.Val) == false {
					unifiedVersion := versionutils.UnifiedVersion(a.Val, padding)
					if i, _ := strconv.Atoi(unifiedVersion[:padding]); i == 4 {
						links[unifiedVersion] = Kernel_webpage + a.Val
					}

					break
				}
			}
		}
	}

	return links
}

func parsePackagePage(url string) (links []string) {
	var reg_deb_all = regexp.MustCompile(`.*_all\.deb`)
	var reg_deb_generic = regexp.MustCompile(`.*generic.*_amd64\.deb`)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken: // End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" && (reg_deb_generic.MatchString(a.Val) || reg_deb_all.MatchString(a.Val)) {
					links = append(links, url+a.Val)

					break
				}
			}
		}
	}
	return
}

func getMostActualKernelVersion(versionsAndLinksMap map[string]string) (version, link string) {
	var keys []string
	for k := range versionsAndLinksMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	version = keys[len(keys)-1]
	link = versionsAndLinksMap[keys[len(keys)-1]]
	return
}

func GetMostActualKernelVersion() (version, link string) {
	links := parseKernelPage()
	version, link = getMostActualKernelVersion(links)
	return
}

func DownloadMostRecentKernelDebs(done chan bool) (version, actual_package_url string) {
	version, actual_package_url = GetMostActualKernelVersion()
	linksToDownload := parsePackagePage(actual_package_url)
	go downloadutils.DownloadFiles(linksToDownload, done)
	return version, actual_package_url
}
