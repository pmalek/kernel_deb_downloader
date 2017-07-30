package ubuntukernelpageutils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pmalek/kernel_deb_downloader/downloadutils"
	"github.com/pmalek/kernel_deb_downloader/versionutils"

	"net/http"
	"regexp"
	"sort"
	"strconv"

	"golang.org/x/net/html"
)

// KernelWebpage - URL pointing to ubuntu's ppa repositorty with Linux kernel's .deb packages
const KernelWebpage = "http://kernel.ubuntu.com/~kernel-ppa/mainline/"

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{} // Use map to record duplicates as we find them.
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == false {
			encountered[elements[v]] = true      // Record this element as an encountered element.
			result = append(result, elements[v]) // Append to result slice.
		}
	}
	return result
}

func parseKernelPage(respBody io.Reader) (links map[string]string) {
	const padding = 2

	z := html.NewTokenizer(respBody)

	links = make(map[string]string, 0) // the same as []string{}

	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		if tt != html.StartTagToken {
			continue
		}

		for _, a := range z.Token().Attr {
			if a.Key != "href" || versionutils.IsAnRCVersion(a.Val) == true {
				continue
			}

			unifiedVersion := versionutils.UnifiedVersion(a.Val, padding)
			if i, _ := strconv.Atoi(unifiedVersion[:padding]); i == 4 {
				links[unifiedVersion] = KernelWebpage + a.Val
			}
		}
	}

	return links
}

func parsePackagePage(url string) (links []string) {
	var regDebAll = regexp.MustCompile(`.*_all\.deb`)
	var regDebGeneric = regexp.MustCompile(`.*generic.*_amd64\.deb`)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could get package webpage %s, received: %v\n", KernelWebpage, err)
		return
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken: // End of the document, we're done
			links = removeDuplicates(links)
			return
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == "href" && (regDebGeneric.MatchString(a.Val) || regDebAll.MatchString(a.Val)) {
					links = append(links, url+a.Val)

					break
				}
			}
		}
	}
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

// GetMostActualKernelVersion returns a pair of strings representing
// version - a canonical kernel version e.g. 040602
// link - a URL where kernel .debs at version @version are stored
func GetMostActualKernelVersion() (version, link string) {
	resp, err := http.Get(KernelWebpage)
	if err != nil {
		fmt.Printf("Could get Ubuntu kernel mainline webpage %s, received: %v\n", KernelWebpage, err)
		return
	}
	defer resp.Body.Close()

	links := parseKernelPage(resp.Body)
	version, link = getMostActualKernelVersion(links)
	return
}

// DownloadKernelDebs downloads Linux kernel .debs from @actualPackageURL
// to the current directory
func DownloadKernelDebs(packageURL string) []string {
	linksToDownload := parsePackagePage(packageURL)
	filenames := downloadutils.DownloadFiles(linksToDownload)
	return filenames
}

// GetChangesFromPackageURL fetches CHANGES file contents from packageURL
// and return a pair of (string) contents of this file and error if not
// successfull
func GetChangesFromPackageURL(packageURL string) (string, error) {
	changesURL := packageURL + "CHANGES"

	response, err := http.Get(changesURL)
	if err != nil {
		return "", err
	} else if response.StatusCode != 200 {
		errStr := fmt.Sprintf("Received %v HTTP status code when downloading CHANGES from %v\n", response.StatusCode, changesURL)
		return "", errors.New(errStr)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(responseData), nil
}
