package downloadutils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func DownloadFile(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %v", err)
		return
	}
	defer output.Close()

	fmt.Printf("Downloading: '%v'...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %v - %v", url, err)
		return
	}
	defer resp.Body.Close()

	n, err := io.Copy(output, resp.Body)
	if err != nil {
		fmt.Printf("Error copying %v - %v", url, err)
		return
	}

	fmt.Printf("Downloaded: %v (%vB)\n", fileName, n)
}

func DownloadFiles(urls []string, done chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(urls)) // Increment the WaitGroup counter.

	for _, url := range urls {
		go func(url string) { // Launch a goroutine to fetch the URL.
			defer wg.Done() // Decrement the counter when the goroutine completes.
			DownloadFile(url)
		}(url)
	}

	wg.Wait() // Wait for all HTTP fetches to complete.
	done <- true
}
