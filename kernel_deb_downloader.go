package main

import (
	"fmt"

	"github.com/pmalek/kernel_deb_downloader/ubuntukernelpageutils"
)

func main() {
	doneDownloading := make(chan bool)
	version, actual_package_url := ubuntukernelpageutils.DownloadMostRecentKernelDebs(doneDownloading)
	fmt.Printf("Most recent (non RC) version: %v, link: %v\n", version, actual_package_url)
	<-doneDownloading
}
