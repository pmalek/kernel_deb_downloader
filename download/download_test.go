package download

import (
	"bytes"
	"math/rand"
	"testing"
	"testing/quick"
	"time"

	"github.com/pmalek/kernel_deb_downloader/http"
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
