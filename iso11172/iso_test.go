// Copyright © 2013-2014 Lawrence E. Bakst. All rights reserved.
package iso11172_test

import . "leb/mpeg-decoder/bitstream"
import . "leb/mpeg-decoder/iso11172"
//import "flag"
import "fmt"
import "testing"

var mpeg1Stream = []byte{
0x00, 0x00, 0x01, 0xB3, 0x02, 0x00, 0x10, 0x14, 0xFF, 0xFF, 0xE0,
0xA0, 0x00, 0x00, 0x01, 0xB8, 0x80, 0x08, 0x00, 0x40, 0x00, 0x00, 0x01,
0x00, 0x00, 0x0F, 0xFF, 0xF8, 0x00, 0x00, 0x01, 0x01, 0xFA, 0x96,
0x52, 0x94, 0x88, 0xAA, 0x25, 0x29, 0x48, 0x88, 0x00, 0x00, 0x01, 0xB7}

func TestShortStream(t *testing.T) {
var ms MpegState

	fmt.Printf("TestShortStream\n")
	bs, _ := NewFromMemory(mpeg1Stream, "r")
	Dump(mpeg1Stream)
	ms.Bitstream = bs
	ms.ReadMPEG1Steam()
}
