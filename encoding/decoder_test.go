package encoding

import (
	"os"
	"testing"
)

var f, _ = os.OpenFile("testdata/silver_open_1.x-ag-binary", os.O_RDONLY, 0600)

func TestDump(t *testing.T) {
	d := NewDecoder(f)

	d.Dump()
}
