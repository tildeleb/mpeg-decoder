package main

import "fmt"
import "flag"
//import "os"
//import "io"
import "leb/mpeg-decoder/bitstream"
import "leb/mpeg-decoder/iso11172"

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
		ms.ReadMPEG1Steam()
	}
}








