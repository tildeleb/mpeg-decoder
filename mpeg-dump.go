package main

//import "fmt"
//import "flag"
//import "os"
//import "io"
import . "leb/mpdm/bitstream"
//import . "leb/mpdm/iso111722"

func main() {

	bs, err := NewFromFile("bike.mpg")
	if err != nil {
		panic("bad filename")
	}

	bs.ReadMPEG1Steam()
}








