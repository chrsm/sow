package main

import (
	"flag"
	"log"
	"os"

	"github.com/chrsm/sow/encoding"
	"github.com/davecgh/go-spew/spew"
)

var (
	fStr = flag.String("file", "", "file to decode (req/resp only, no headers plz kthx)")
)

func main() {
	flag.Parse()
	assert(*fStr != "", "need path to file that has req/resp data (-file)")

	f, err := os.Open(*fStr)
	assert(err == nil, "Couldn't open file: %s", err)
	defer f.Close()

	dec := encoding.NewDecoder(f)
	v, err := dec.DecodeRaw()
	assert(err == nil, "Couldn't decode file: %s", err)

	log.Println("decoded data: big dump of crap incoming")
	spew.Dump(v)
}

func assert(cond bool, fmts string, vars ...interface{}) {
	if !cond {
		log.Fatalf(fmts, vars...)
	}
}
