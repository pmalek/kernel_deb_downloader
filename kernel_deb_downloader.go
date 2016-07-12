package main

import (
	"flag"
	"fmt"

	"github.com/pmalek/kernel_deb_downloader/ubuntukernelpageutils"
)

func main() {
	onlyPrintVersion := flag.Bool("n", false, "Print newest version - do not download the .debs")
	flag.Parse()

	if *onlyPrintVersion == false {
		doneDownloading := make(chan bool)
		version, actualPackageURL := ubuntukernelpageutils.DownloadMostRecentKernelDebs(doneDownloading)
		fmt.Printf("Most recent (non RC) version: %v, link: %v\n", version, actualPackageURL)
		<-doneDownloading
	} else {
		version, actualPackageURL := ubuntukernelpageutils.GetMostActualKernelVersion()
		fmt.Printf("Most recent (non RC) version: %v, link: %v\n", version, actualPackageURL)
	}

}
