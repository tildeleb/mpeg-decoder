// Copyright © 2003 and 2014 Lawrence E. Bakst. All rights reserved.
// THIS SOURCE CODE IS THE PROPRIETARY INTELLECTUAL PROPERTY AND CONFIDENTIAL
// INFORMATION OF LAWRENCE E. BAKST AND IS PROTECTED UNDER U.S. AND
// INTERNATIONAL LAW. ANY USE OF THIS SOURCE CODE WITHOUT THE PRIOR WRITTEN
// AUTHORIZATION OF LAWRENCE E. BAKST IS STRICTLY PROHIBITED.

package main

import "fmt"
import "flag"
//import "os"
//import "io"
import "leb/mpeg-decoder/bitstream"
import "leb/mpeg-decoder/iso11172"

var phdf = flag.Bool("phd", false, "print headers")
var pvsf = flag.Bool("pvs", false, "print video slices")
var pmbf = flag.Bool("pmb", false, "print macro blocks")
var pbcf = flag.Bool("pbc", false, "print block coefficients")
var prmbf = flag.Bool("prmb", false, "print raw macro blocks")
var rmbf = flag.Bool("rmb", true, "read macro blocks")
var vf = flag.Bool("v", false, "verbose; turns on all printing")
var psf = flag.Bool("ps", false, "print stats")

var from = flag.Int("from", 0, "start at frame #")
var to = flag.Int("to", 9999999, "stop at frame #")

func chk(i iso11172.Mpeg1Parser) {
	_, ok := i.(iso11172.Mpeg1Parser)
	//fmt.Printf("%T | %T\n", i, x)
	if ok {
		//fmt.Printf("OK")
	} else {
		fmt.Printf("BAD")	
	}
}

func main() {
	flag.Parse()
	if *vf {
		*phdf, *pvsf, *pmbf = true, true, true
	}
	ms := iso11172.New(*phdf, *pvsf, *pmbf, *pbcf, *prmbf, *rmbf)
	chk(ms)
	ms.CBPS = make(map[string]int, 100)

	for i := 0; i < flag.NArg(); i++ {
		//fmt.Printf("arg %d=|%s|\n", i, flag.Arg(i))
		bs, err := bitstream.NewFromFile(flag.Arg(i), "r")
		if err != nil {
			panic("bad filename")
		}
		ms.Bitstream = bs
		ms.ReadMPEG1Steam(*from, *to)
		//fmt.Printf("ms.MpegStats=%#v\n", ms.MpegStats)
		if *psf {
			ms.PrintStats()
		}
	}
}








