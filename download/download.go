package download

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pmalek/kernel_deb_downloader/http"
	"github.com/pmalek/pb"
)

// ToWriter downloads contents from url using http client
// and writes it into the io.Writer - out.
// It returns number of bytes written and an error
func ToWriter(client http.Getter, out io.Writer, url string, progressBar ...*pb.ProgressBar) (int64, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error downloading %v, error : %v", url, err)
	}
	defer resp.Body.Close()

	if len(progressBar) > 1 {
		return 0, fmt.Errorf("Passed in more than 1 progressBar")
	}

	var reader io.Reader = resp.Body
	for _, p := range progressBar {
		reader = p.NewProxyReader(resp.Body)
	}

	var n int64
	if n, err = io.Copy(out, reader); err != nil {
		return 0, fmt.Errorf("error copying to %v, error : %v", fileName, err)
	}

	return n, nil
}

func httpFileSizeWithHEAD(client http.Header, url string) (int64, error) {
	resp, err := client.Head(url)
	if err != nil {
		return -1, fmt.Errorf("error HEADing %v, error : %v", url, err)
	}
	defer resp.Body.Close()

	return resp.ContentLength, nil
}

func fileNameFromURL(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

// ToFiles downloads all the files from @urls in package
// and puts the in the current directory
func ToFiles(client http.GetterHeader, urls []string) []string {
	filenames := make([]string, 0, len(urls))

	pool, err := pb.StartPool()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(urls)) // Increment the WaitGroup counter.

	for _, url := range urls {
		go func(url string) { // Launch a goroutine to fetch the URL.
			defer wg.Done() // Decrement the counter when the goroutine completes.

			fileSize, err := httpFileSizeWithHEAD(client, url)
			if err != nil {
				log.Fatal(err)
				return
			} else if fileSize < 0 {
				log.Fatalf("Requesting HEAD on %s return %d ContentLength", url, fileSize)
				return
			}

			fileName := fileNameFromURL(url)
			progressBar := pb.
				New64(fileSize).
				SetUnits(pb.U_BYTES).
				Prefix(fmt.Sprintf("%-76s", fileName))
			progressBar.ShowSpeed = true
			pool.Add(progressBar)

			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Error creating file %v, error : %v\n", fileName, err)
				return
			}
			defer file.Close()

			if _, err = ToWriter(client, file, url, progressBar); err != nil {
				filenames = append(filenames, fileName)
			}
		}(url)
	}

	defer pool.Stop()
	wg.Wait() // Wait for all HTTP fetches to complete.

	return filenames
}
