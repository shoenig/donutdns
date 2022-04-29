package extract

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/shoenig/test/must"
)

func openSample(t *testing.T, filename string) io.Reader {
	b, err := ioutil.ReadFile(fmt.Sprintf("samples/%s", filename))
	must.NoError(t, err)
	return bytes.NewBuffer(b)
}

func TestExtractor_Extract(t *testing.T) {

	try := func(filename, re string, exp int) {
		ex := New(Generic)
		result, err := ex.Extract(openSample(t, filename))
		must.NoError(t, err)
		must.EqOp(t, exp, result.Len())
	}

	try("KADhosts.txt", Generic, 6)
	try("w3kbl.txt", Generic, 7)
	try("adaway.txt", Generic, 3)
	try("Admiral.txt", Generic, 1)
	try("adservers.txt", Generic, 2)
	try("hostsVN.txt", Generic, 3)
	try("AntiMalwareHosts.txt", Generic, 6) // matches some ipv4 addresses
	try("Prigent-Crypto.txt", Generic, 4)
	try("notrack-malware.txt", Generic, 4)
	try("abuse.txt", Generic, 5)
}
