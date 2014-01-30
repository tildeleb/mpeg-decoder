// Copyright © 2012-2013 Lawrence E. Bakst. All rights reserved.

package iso11172_test

import "leb/mpeg-decoder/bitstream"
import "leb/mpeg-decoder/iso11172"
//import "flag"
import "fmt"
import "math/rand"
import "testing"

var r = rand.Float64

type DCTvlc struct {
	r		int
	l		int
	d		uint32
	vlc		uint32		// converted from the above
	v		*[]byte		// warning only the low order nibble of each byte is used to form the VLC. It was easier to type in that way.
}

var vlcmap map[uint32]DCTvlc = make(map[uint32]DCTvlc, 10)

var vlcs []DCTvlc =  []DCTvlc{
//	{0, 1, 2, 0, &[]byte{0x1}},
	{0, 1, 3, 0, &[]byte{0x3}},
	{0, 2, 5, 0, &[]byte{0x4}},
	{0, 3, 6, 0, &[]byte{0x2, 0x1}},
	{0, 4, 8, 0, &[]byte{0x0, 0x6}},
	{0, 5, 9, 0, &[]byte{0x2, 0x6}},
	{0, 6, 9, 0, &[]byte{0x2, 0x1}},
	{0, 7, 11, 0, &[]byte{0x0, 0x2, 0x2}},
	{0, 8, 13, 0, &[]byte{0x0, 0x1, 0xD}},
	{0, 9, 13, 0, &[]byte{0x0, 0x1, 0x8}},
	{0, 10, 13, 0, &[]byte{0x0, 0x1, 0x3}},
	{0, 11, 13, 0, &[]byte{0x0, 0x1, 0x0}},
	{0, 12, 14, 0, &[]byte{0x0, 0x0, 0xD, 0x0}},
	{0, 13, 14, 0, &[]byte{0x0, 0x0, 0xC, 0x1}},
	{0, 14, 14, 0, &[]byte{0x0, 0x0, 0xC, 0x0}},
	{0, 15, 14, 0, &[]byte{0x0, 0x0, 0xB, 0x1}},

	{0, 16, 15, 0, &[]byte{0x0, 0x0, 0x7, 0x3}},
	{0, 17, 15, 0, &[]byte{0x0, 0x0, 0x7, 0x2}},
	{0, 18, 15, 0, &[]byte{0x0, 0x0, 0x7, 0x1}},
	{0, 19, 15, 0, &[]byte{0x0, 0x0, 0x7, 0x0}},
	{0, 20, 15, 0, &[]byte{0x0, 0x0, 0x6, 0x3}},
	{0, 21, 15, 0, &[]byte{0x0, 0x0, 0x6, 0x2}},
	{0, 22, 15, 0, &[]byte{0x0, 0x0, 0x6, 0x1}},
	{0, 23, 15, 0, &[]byte{0x0, 0x0, 0x6, 0x0}},

	{0, 24, 15, 0, &[]byte{0x0, 0x0, 0x5, 0x3}},
	{0, 25, 15, 0, &[]byte{0x0, 0x0, 0x5, 0x2}},
	{0, 26, 15, 0, &[]byte{0x0, 0x0, 0x5, 0x1}},
	{0, 27, 15, 0, &[]byte{0x0, 0x0, 0x5, 0x0}},
	{0, 28, 15, 0, &[]byte{0x0, 0x0, 0x4, 0x3}},
	{0, 29, 15, 0, &[]byte{0x0, 0x0, 0x4, 0x2}},
	{0, 30, 15, 0, &[]byte{0x0, 0x0, 0x4, 0x1}},
	{0, 31, 15, 0, &[]byte{0x0, 0x0, 0x4, 0x0}},

	{0, 32, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x0}},
	{0, 33, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x7}},
	{0, 34, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x6}},
	{0, 35, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x5}},
	{0, 36, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x4}},
	{0, 37, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x3}},
	{0, 38, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x2}},
	{0, 39, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x1}},
	{0, 40, 16, 0, &[]byte{0x0, 0x0, 0x2, 0x0}},


	{1, 1, 4, 0, &[]byte{0x3}},
	{1, 2, 7, 0, &[]byte{0x1, 0x2}},
	{1, 3, 9, 0, &[]byte{0x2, 0x5}},
	{1, 4, 11, 0, &[]byte{0x0, 0x3, 0x0}},
	{1, 5, 13, 0, &[]byte{0x0, 0x1, 0xB}},
	{1, 6, 14, 0, &[]byte{0x0, 0x0, 0xB, 0x0}},
	{1, 7, 14, 0, &[]byte{0x0, 0x0, 0xA, 0x1}},

	{1, 8, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x7}},
	{1, 9, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x6}},
	{1, 10, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x5}},
	{1, 11, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x4}},
	{1, 12, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x3}},
	{1, 13, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x2}},
	{1, 14, 16, 0, &[]byte{0x0, 0x0, 0x3, 0x1}},

	{1, 15, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x3}},
	{1, 16, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x2}},
	{1, 17, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x1}},
	{1, 18, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x0}},


	{2, 1, 5, 0, &[]byte{0x5}},
	{2, 2, 8, 0, &[]byte{0x0, 0x4}},
	{2, 3, 11, 0, &[]byte{0x0, 0x2, 0x3}},
	{2, 4, 13, 0, &[]byte{0x0, 0x1, 0x4}},
	{2, 5, 14, 0, &[]byte{0x0, 0x0, 0xA, 0x0}},

	{3, 1, 6, 0, &[]byte{0x3, 0x1}},
	{3, 2, 9, 0, &[]byte{0x2, 0x4}},
	{3, 3, 13, 0, &[]byte{0x0, 0x1, 0xC}},
	{3, 4, 14, 0, &[]byte{0x0, 0x0, 0x9, 0x1}},

	{4, 1, 6, 0, &[]byte{0x3, 0x0}},
	{4, 2, 11, 0, &[]byte{0x0, 0x3, 0x3}},
	{4, 3, 13, 0, &[]byte{0x0, 0x1, 0x2}},

	{5, 1, 7, 0, &[]byte{0x1, 0x3}},
	{5, 2, 11, 0, &[]byte{0x0, 0x2, 0x1}},
	{5, 3, 14, 0, &[]byte{0x0, 0x0, 0x9, 0x0}},

	{6, 1, 7, 0, &[]byte{0x1, 0x1}},
	{6, 2, 13, 0, &[]byte{0x0, 0x1, 0xE}},
	{6, 3, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x4}},

	{7, 1, 7, 0, &[]byte{0x1, 0x0}},
	{7, 2, 13, 0, &[]byte{0x0, 0x1, 0x5}},

	{8, 1, 8, 0, &[]byte{0x0, 0x7}},
	{8, 2, 13, 0, &[]byte{0x0, 0x1, 0x1}},

	{9, 1, 8, 0, &[]byte{0x0, 0x5}},
	{9, 2, 14, 0, &[]byte{0x0, 0x0, 0x8, 0x1}},

	{10, 1, 9, 0, &[]byte{0x2, 0x7}},
	{10, 2, 14, 0, &[]byte{0x0, 0x0, 0x8, 0x0}},

	{11, 1, 9, 0, &[]byte{0x2, 0x3}},
	{11, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xA}},

	{12, 1, 9, 0, &[]byte{0x2, 0x2}},
	{12, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x9}},

	{13, 1, 9, 0, &[]byte{0x2, 0x0}},
	{13, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x8}},

	{14, 1, 11, 0, &[]byte{0x0, 0x3, 0x2}},
	{14, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x7}},

	{15, 1, 11, 0, &[]byte{0x0, 0x3, 0x1}},
	{15, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x6}},

	{16, 1, 11, 0, &[]byte{0x0, 0x2, 0x0}},
	{16, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0x5}},

	{17, 1, 13, 0, &[]byte{0x0, 0x1, 0xF}},
	{18, 1, 13, 0, &[]byte{0x0, 0x1, 0xA}},
	{19, 1, 13, 0, &[]byte{0x0, 0x1, 0x9}},
	{20, 1, 13, 0, &[]byte{0x0, 0x1, 0x7}},
	{21, 1, 13, 0, &[]byte{0x0, 0x1, 0x6}},

	{22, 1, 14, 0, &[]byte{0x0, 0x0, 0xF, 0x1}},
	{23, 1, 14, 0, &[]byte{0x0, 0x0, 0xF, 0x0}},
	{24, 1, 14, 0, &[]byte{0x0, 0x0, 0xE, 0x1}},
	{25, 1, 14, 0, &[]byte{0x0, 0x0, 0xE, 0x0}},
	{26, 1, 14, 0, &[]byte{0x0, 0x0, 0xD, 0x1}},

	{27, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xF}},
	{28, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xE}},
	{29, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xD}},
	{30, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xC}},
	{31, 2, 17, 0, &[]byte{0x0, 0x0, 0x1, 0xB}},
}

func convert(b *[]byte, bits uint32, sign uint32, signif bool) (vlc uint32) {
	if signif {
		vlc = 1 // make leading zeros significant
	}
	bits--
	for _, v := range *b {
		if bits > 4 {
			vlc <<= 4
			vlc |= uint32(v)
			bits -= 4
		} else {
			vlc <<= uint(bits)+1
			vlc |= (uint32(v)<<1)|sign
			return
		}
	}
	panic("convert: not enough bytes")
}

func checkDup(r, l int) {
	for _, v := range vlcmap {
		if v.r == r && v.l == l {
			panic("dup r/l")
		}
	}
}


func TestVlc(t *testing.T) {
//var ms MpegState

	for k, v := range vlcs {
		vlcPlus := convert(v.v, v.d, 0, true)
		vlcMinus := convert(v.v, v.d, 1, true)
		fmt.Printf("k=%d, v.r=%d, v.l=%d, v.d=%d, vlcPlus=0x%x, vlcMinus=0x%x\n", k, v.r, v.l, v.d, vlcPlus, vlcMinus)
		_, ok := vlcmap[vlcPlus]
		if ok {
			fmt.Printf("TestVlc: 1 - map entry already used, k=%d, vlc=0x%x\n", k, vlcPlus)
		}
		checkDup(v.r, v.l)
		checkDup(v.r, -v.l)
		v.vlc = vlcPlus
		vlcmap[vlcPlus] = v

		_, ok = vlcmap[vlcMinus]
		if ok {
			fmt.Printf("TestVlc: 2 - map entry already used, k=%d, vlc=0x%x\n", k, vlcMinus)
		}
		v.vlc = vlcMinus
		v.l = -v.l
		vlcmap[vlcMinus] = v
		//vlcmap[vlc].vlc = vlc
	}
}

func rbetween(a int, b int) int {
	rf := r()
	diff := float64(b - a + 1)
	r2 := rf * diff
	r3 := r2 + float64(a)
	//fmt.Printf("rbetween: a=%d, b=%d, rf=%f, diff=%f, r2=%f, r3=%f\n", a, b, rf, diff, r2, r3)
	ret := int(r3)
	return ret
}


func TestVlc2(t *testing.T) {
const maxbits uint64 = 1280000000
var tdata	uint64	// buffer where bits are stored, drained to less than 8 bits, after data is added
var tblen	uint64	// used to keep track of how many bit used in the above buffer
var tbitlen uint64 // total bits used
var data []byte
var cnt int
var zero uint32
var ms iso11172.MpegState
var sgn = func(s int) int {
	if (s == 1) {
		return -1
	} else {
		return 1
	}
}

	fmt.Printf("TestVlc2\n")
	tableLen := len(vlcs)

	for tbitlen < maxbits {
		i := rbetween(0, tableLen-1)
		s := rbetween(0, 1)
		r, l, b, d := vlcs[i].r, vlcs[i].l, vlcs[i].v,  vlcs[i].d
		// convert(b *[]byte, bits uint32, sign uint32, signif bool) (vlc uint32)
		vlc := convert(b, vlcs[i].d, uint32(s), false)
		// func Put(bits []byte, value uint32, blen uint, tdata *uint64, tblen *uint64, tbitlen *uint64) {
		//fmt.Printf("TestVlc2: cnt=%d, i=%d, s=%d, r=%d, l=%d, vlc=0x%x/%d\n", cnt, i, s, r, l*sgn(s), vlc, d)
		r++
		l++
		l = sgn(r)
		bitstream.Put(&data, vlc, uint(d), &tdata, &tblen, &tbitlen)
		cnt++
	}
	bitstream.Put(&data, iso11172.EOB, 2, &tdata, &tblen, &tbitlen)
	bitstream.Put(&data, zero, 32, &tdata, &tblen, &tbitlen)
	bitstream.Put(&data, zero, 32, &tdata, &tblen, &tbitlen)
	fmt.Printf("TestVlc2: len(data)=%d, %d bits in buffer\n", len(data), tbitlen)
	//bitstream.Dump(data, len(data))
	ms.Bitstream, _ = bitstream.NewFromMemory(data, "r")
	for ms.Peekbits(2) != iso11172.EOB {
		r, l := ms.DecodeDCTCoeff(false)
		//fmt.Printf("TestVlc2: r=%d, l=%d\n", r, l)
		r++
		l++
	}
}

