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

func parsePackagePage(respBody io.Reader, packageURL string) (links []string) {
	var regDebAll = regexp.MustCompile(`.*_all\.deb`)
	var regDebGeneric = regexp.MustCompile(`.*generic.*_amd64\.deb`)

	z := html.NewTokenizer(respBody)

	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {

		if tt != html.StartTagToken {
			continue
		}

		t := z.Token()

		for _, a := range t.Attr {
			if a.Key == "href" && (regDebGeneric.MatchString(a.Val) || regDebAll.MatchString(a.Val)) {
				links = append(links, packageURL+a.Val)
				break
			}
		}
	}

	links = removeDuplicates(links)
	return
}

func getMostActualKernelVersion(versionsAndLinksMap map[string]string) (version, link string) {
	if len(versionsAndLinksMap) == 0 {
		return
	}

	var keys []string
	for k := range versionsAndLinksMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	version = keys[len(keys)-1]
	link = versionsAndLinksMap[keys[len(keys)-1]]

	return
}

type httpGetter interface {
	Get(string) (*http.Response, error)
}

// GetMostActualKernelVersion returns a pair of strings representing
// version - a canonical kernel version e.g. 040602
// link - a URL where kernel .debs at version @version are stored
func GetMostActualKernelVersion(client httpGetter) (version, link string) {
	resp, err := client.Get(KernelWebpage)

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
func DownloadKernelDebs(packageURL string) ([]string, error) {
	resp, err := http.Get(packageURL)
	if err != nil {
		fmt.Printf("Could get package webpage %s, received: %v\n", KernelWebpage, err)
		return nil, err
	}
	defer resp.Body.Close()

	linksToDownload := parsePackagePage(resp.Body, packageURL)
	filenames := downloadutils.DownloadFiles(linksToDownload)
	return filenames, nil
}

// GetChangesFromPackageURL fetches CHANGES file contents from packageURL
// and returns a pair of (string) contents of this file and an error
// if not successful
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
