package ubuntukernelpageutils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pmalek/kernel_deb_downloader/download"
	"github.com/pmalek/kernel_deb_downloader/http"
	"github.com/pmalek/kernel_deb_downloader/versionutils"

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

	if len(keys) < 1 {
		return
	}

	version = keys[len(keys)-1]
	link = versionsAndLinksMap[keys[len(keys)-1]]

	return
}

func getMostActualKernelVersionFromMajorVersion(desiredVersion string, versionsAndLinksMap map[string]string) (version, link string) {
	if len(versionsAndLinksMap) == 0 {
		return
	}

	versionLen := 0

	switch strings.Count(desiredVersion, ".") {
	case 2:
		desiredVersion = versionutils.UnifiedVersion(desiredVersion, 2)[:6]
		versionLen = 6
	case 1:
		desiredVersion = versionutils.UnifiedVersion(desiredVersion, 2)[:4]
		versionLen = 4
	case 0:
		desiredVersion = versionutils.UnifiedVersion(desiredVersion, 2)[:2]
		versionLen = 2
	default:
		return
	}

	var keys []string
	for k := range versionsAndLinksMap {
		if len(k) < versionLen {
			continue
		}

		version := k[:versionLen]
		if strings.Compare(version, desiredVersion) == 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	if len(keys) < 1 {
		return
	}

	version = keys[len(keys)-1]
	link = versionsAndLinksMap[keys[len(keys)-1]]

	return
}

// GetMostActualKernelVersion returns a pair of strings representing
// version - a canonical kernel version e.g. 040602
// link - a URL where kernel .debs at version @version are stored
func GetMostActualKernelVersion(client http.Getter) (version, link string, err error) {
	resp, err := client.Get(KernelWebpage)

	if err != nil {
		return "", "",
			fmt.Errorf("Could get Ubuntu kernel mainline webpage %s, received error: %v", KernelWebpage, err)
	}
	defer resp.Body.Close()

	links := parseKernelPage(resp.Body)
	version, link = getMostActualKernelVersion(links)
	return version, link, nil
}

// GetMostActualKernelVersionFromMajorVersion returns a pair of strings representing
// version - a canonical kernel version e.g. 040602
// link - a URL where kernel .debs at version @version are stored
// from major version release branch e.g. 4.13
func GetMostActualKernelVersionFromDesiredVersion(desiredVersion string, client http.Getter) (version, link string, err error) {
	resp, err := client.Get(KernelWebpage)

	if err != nil {
		return "", "",
			fmt.Errorf("Could get Ubuntu kernel mainline webpage %s, received error: %v", KernelWebpage, err)
	}
	defer resp.Body.Close()

	links := parseKernelPage(resp.Body)
	version, link = getMostActualKernelVersionFromMajorVersion(desiredVersion, links)
	return version, link, nil
}

// DownloadKernelDebs downloads Linux kernel .debs from @actualPackageURL
// to the current directory
func DownloadKernelDebs(client http.GetterHeader, packageURL string) ([]string, error) {
	resp, err := client.Get(packageURL)
	if err != nil {
		fmt.Printf("Could get package webpage %s, received: %v\n", KernelWebpage, err)
		return nil, err
	}
	defer resp.Body.Close()

	linksToDownload := parsePackagePage(resp.Body, packageURL)
	filenames := download.ToFiles(client, linksToDownload)
	return filenames, nil
}

// GetChangesFromPackageURL fetches CHANGES file contents from packageURL
// and returns contents of this file and an error if not successful
func GetChangesFromPackageURL(client http.Getter, packageURL string) (string, error) {
	changesURL := packageURL + "CHANGES"

	response, err := client.Get(changesURL)
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
