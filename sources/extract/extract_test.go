package extract

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/shoenig/test/must"
)

func openSample(t *testing.T, filename string) io.Reader {
	b, err := os.ReadFile(fmt.Sprintf("samples/%s", filename))
	must.NoError(t, err)
	return bytes.NewBuffer(b)
}

func TestExtractor_Extract(t *testing.T) {
	cases := []struct {
		file string
		mode string
		exp  int
	}{
		{"KADhosts.txt", Generic, 6},
		{"w3kbl.txt", Generic, 7},
		{"adaway.txt", Generic, 3},
		{"Admiral.txt", Generic, 1},
		{"adservers.txt", Generic, 2},
		{"hostsVN.txt", Generic, 3},
		{"AntiMalwareHosts.txt", Generic, 6}, // matches some ipv4 addresses
		{"Prigent-Crypto.txt", Generic, 4},
		{"notrack-malware.txt", Generic, 4},
		{"abuse.txt", Generic, 5},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			ex := New(tc.mode)
			result, err := ex.Extract(openSample(t, tc.file))
			must.NoError(t, err)
			must.Eq(t, tc.exp, result.Size(), must.Sprintf("result: %v", result))
		})
	}
}
