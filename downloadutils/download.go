package downloadutils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/pmalek/pb"
)

func httpFileSizeResp(resp *http.Response) (int64, error) {
	fileSizeStr := resp.Header.Get("Content-Length")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return fileSize, nil
}

func httpFileSizeWithHEAD(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Printf("Error HEADing %v, error : %v\n", url, err)
		return -1, err
	}

	defer resp.Body.Close()

	fileSizeStr := resp.Header.Get("Content-Length")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return fileSize, nil
}

func fileNameFromURL(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

// DownloadFile downloads a file from under the URL @url
// and puts it in the current directory under the same name
// (pretty much the same as wget http://someurl.com/file.ext )
func DownloadFile(url string, progressBar *pb.ProgressBar) (string, error) {
	fileName := fileNameFromURL(url)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %v, error : %v\n", fileName, err)
		return "", err
	}
	defer output.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %v, error : %v\n", url, err)
		return "", err
	}
	defer resp.Body.Close()

	proxyRespReader := progressBar.NewProxyReader(resp.Body)

	//var n int64
	if _, err = io.Copy(output, proxyRespReader); err != nil {
		fmt.Printf("Error copying to %v, error : %v\n", fileName, err)
		return "", err
	}

	//fmt.Printf("Downloaded: %v (%vB)\n", fileName, n)
	return fileName, nil
}

// DownloadFiles downloads all the files from @urls in package
// and puts the in the current directory
func DownloadFiles(urls []string) []string {
	filenames := make([]string, 0, len(urls))

	pool, err := pb.StartPool()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			fileSize, err := httpFileSizeWithHEAD(url)
			if err != nil {
				log.Fatal(err)
				return
			}

			fileName := fileNameFromURL(url)
			progressBar := pb.New64(fileSize).SetUnits(pb.U_BYTES)
			progressBar.Prefix(fmt.Sprintf("%-76s", fileName))
			progressBar.ShowSpeed = true
			pool.Add(progressBar)

			if filename, err := DownloadFile(url, progressBar); err != nil {
				filenames = append(filenames, filename)
			}
		}(url)
	}
	defer pool.Stop()
	wg.Wait()

	return filenames
}
