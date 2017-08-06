package download

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/pmalek/kernel_deb_downloader/http"
	"github.com/pmalek/pb"
)

func Test_ToWriter_quick(t *testing.T) {
	f := func(httpContent string) bool {
		client := http.MockedClient{}
		client.SetResponse(httpContent)

		buff := &bytes.Buffer{}
		n, err := ToWriter(client, buff, "")

		if n != int64(buff.Len()) {
			t.Errorf("Wrong number of bytes returned: %d, expected %d", n, buff.Len())
			return false
		}

		if err != nil {
			t.Errorf("No error expected yet received %q", err)
			return false
		}

		if buff.String() != httpContent {
			t.Errorf("Wrong content returned: %q, expected: %q", buff.String(), httpContent)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		MaxCount: 10000,
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}

func Test_ToWriter_Error(t *testing.T) {
	client := http.MockedClient{}
	client.SetError(fmt.Errorf("Some error"))

	buff := &bytes.Buffer{}
	_, err := ToWriter(client, buff, "")

	if err == nil {
		t.Errorf("Error expected yet received nil error")
	}
}

func Test_ToWriter_OneProgressBar(t *testing.T) {
	expectedContent := "Content"
	client := http.MockedClient{}
	client.SetResponse(expectedContent)

	buff := &bytes.Buffer{}
	_, err := ToWriter(client, buff, "", pb.New(0))

	if err != nil {
		t.Errorf("Error returned yet not expected, returned %q", err)
	}

	if buff.String() != expectedContent {
		t.Errorf("Expected %q but returned %q", expectedContent, buff.String())
	}
}

func Test_ToWriter_Error_ToManyProgressBars(t *testing.T) {
	client := http.MockedClient{}

	buff := &bytes.Buffer{}
	_, err := ToWriter(client, buff, "", pb.New(0), pb.New(0))

	if err == nil {
		t.Errorf("Error expected yet received nil error")
	}
}

func Test_httpFileSizeWithHead(t *testing.T) {
	content := "Content"

	client := http.MockedClient{}
	client.SetResponse(content)

	n, err := httpFileSizeWithHEAD(client, "")

	if err != nil {
		t.Errorf("Error returned yet not expected, returned %q", err)
	}
	if n != int64(len(content)) {
		t.Errorf("Expected %d yet returned %d", len(content), n)
	}
}

func Test_httpFileSizeWithHead_Error(t *testing.T) {
	client := http.MockedClient{}
	client.SetError(errors.New(""))

	_, err := httpFileSizeWithHEAD(client, "")

	if err == nil {
		t.Errorf("Error expected yet received nil error")
	}
}

func Test_fileNameFromURL(t *testing.T) {
	f := func(url string) bool {
		fileName := fileNameFromURL(url)

		if strings.Contains(fileName, "/") {
			t.Errorf("Extracted fileName %q contains / yet it shouldn't ", fileName)
			return false
		}

		if len(fileName) > len(url) {
			t.Errorf("Extracted fileName %q is longer than URL from which it was extracted %q "+
				"yet it shouldn't ", fileName, url)
			return false
		}

		return true
	}

	cfg := &quick.Config{
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		MaxCount: 10000,
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}
