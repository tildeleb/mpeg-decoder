package main

import "fmt"
import "flag"
//import "os"
//import "io"
import "leb/mpeg-decoder/bitstream"
import "leb/mpeg-decoder/iso11172"

var rmbf = flag.Bool("rmb", true, "read macro blocks")
var pmbf = flag.Bool("pmb", false, "print macro blocks")
var from = flag.Int("from", 0, "start at frame #")
var to = flag.Int("to", 99999, "stop at frame #")

func main() {
	var ms iso11172.MpegState

	flag.Parse()
	for i := 0; i < flag.NArg(); i++ {
		fmt.Printf("arg %d=|%s|\n", i, flag.Arg(i))
		bs, err := bitstream.NewFromFile(flag.Arg(i), "r")
		if err != nil {
			panic("bad filename")
		}
		ms.Bitstream = bs
		ms.ReadMPEG1Steam(*from, *to, *rmbf, *pmbf)
	}
}








