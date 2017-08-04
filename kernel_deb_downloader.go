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

	version, packageURL := ubuntukernelpageutils.GetMostActualKernelVersion(http.DefaultClient)
	fmt.Printf("Most recent (non RC) version: %v, link: %v\n", version, packageURL)

	if showChanges {
		if changes, err := ubuntukernelpageutils.GetChangesFromPackageURL(packageURL); err != nil {
			fmt.Printf("Error downloading changes: %v", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("Changes: \n%v", changes)
		}
	}

	if onlyPrintVersion == false {
		ubuntukernelpageutils.DownloadKernelDebs(packageURL)
	}

}
