package download

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/pmalek/kernel_deb_downloader/http"
)

// ToWriter downloads contents from url using http client
// and writes it into the io.Writer - out.
// It returns number of bytes written and an error
func ToWriter(client http.Getter, out io.Writer, url string) (int64, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error downloading %v, error : %v", url, err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error copying to %v, error : %v", fileName, err)
	}

	return n, nil
}

// ToFile downloads contents from url using http client
// and writes it to a file named the same as in the source URL
// It returns number of bytes written and an error
func ToFile(client http.Getter, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %v, error : %v\n", fileName, err)
		return "", err
	}
	defer output.Close()

	fmt.Printf("Downloading: '%v'...\n", url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %v, error : %v\n", url, err)
		return "", err
	}
	defer resp.Body.Close()

	n, copyErr := io.Copy(output, resp.Body)
	if copyErr != nil {
		fmt.Printf("Error copying to %v, error : %v\n", fileName, err)
		return "", err
	}

	fmt.Printf("Downloaded: %v (%vB)\n", fileName, n)
	return fileName, nil
}

// ToFiles downloads all the files from @urls in package
// and puts the in the current directory
func ToFiles(client http.Getter, urls []string) []string {
	filenames := make([]string, 0, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls)) // Increment the WaitGroup counter.

	for _, url := range urls {
		go func(url string) { // Launch a goroutine to fetch the URL.
			defer wg.Done() // Decrement the counter when the goroutine completes.

			tokens := strings.Split(url, "/")
			fileName := tokens[len(tokens)-1]
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Error creating file %v, error : %v\n", fileName, err)
				return
			}
			defer file.Close()

			if _, err := ToWriter(client, file, url); err != nil {
				filenames = append(filenames, fileName)
			}
		}(url)
	}

	wg.Wait() // Wait for all HTTP fetches to complete.

	return filenames
}
