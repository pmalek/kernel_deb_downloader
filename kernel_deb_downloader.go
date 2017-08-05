package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/pmalek/kernel_deb_downloader/ubuntukernelpageutils"
)

var (
	onlyPrintVersion bool
	showChanges      bool
)

func init() {
	flag.BoolVar(&onlyPrintVersion, "n", false, "Print newest version - do not download the .debs")
	flag.BoolVar(&showChanges, "c", false, "Show changes included in particular kernel package")
}

func main() {
	flag.Parse()

	version, packageURL, err := ubuntukernelpageutils.GetMostActualKernelVersion(http.DefaultClient)
	if err != nil {
		fmt.Printf("Error connecting to Ubuntu's kernel ppa webpage, error: %q", err)
		os.Exit(1)
	}

	fmt.Printf("Most recent (non RC) version: %v, link: %v\n", version, packageURL)

	if showChanges {
		if changes, err := ubuntukernelpageutils.GetChangesFromPackageURL(http.DefaultClient, packageURL); err != nil {
			fmt.Printf("Error downloading changes: %v", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("Changes: \n%v", changes)
		}
	}

	if onlyPrintVersion == false {
		_, err = ubuntukernelpageutils.DownloadKernelDebs(http.DefaultClient, packageURL)
		if err != nil {
			fmt.Printf("Error downloading .deb files: %q", err)
		}
	}

}
