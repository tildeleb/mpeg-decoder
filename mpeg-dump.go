package main

//import "fmt"
//import "flag"
//import "os"
//import "io"
import "leb/mpeg-decoder/bitstream"
import "leb/mpeg-decoder/iso11172"

func main() {
var ms iso11172.MpegState

	bs, err := bitstream.NewFromFile("bike.mpg")
	if err != nil {
		panic("bad filename")
	}
	ms.Bitstream = bs
	ms.ReadMPEG1Steam()
}








