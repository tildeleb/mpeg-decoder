// Copyright © 2012-2013 Lawrence E. Bakst. All rights reserved.

package bitstream_test

import . "leb/mpeg-decoder/bitstream"
//import "flag"
import "fmt"
import "math/rand"
import "math"
import "testing"

var r = rand.Float64
const maxbitlen = 32
const nBytes = 1000000 // 500*1024*1024
const nEntries = 1000 * 1000 /// 
var bits []byte
type info struct {
	blen	uint
	value	uint32
}
var data []info
var dist [33]uint


func rbetween(a int, b int) int {
	rf := r()
	diff := float64(b - a + 1)
	r2 := rf * diff
	r3 := r2 + float64(a)
	//fmt.Printf("rbetween: a=%d, b=%d, rf=%f, diff=%f, r2=%f, r3=%f\n", a, b, rf, diff, r2, r3)
	ret := int(r3)
	return ret
}

func dumpDist() {
	for k, v := range dist {
		fmt.Printf("dist[%d]=%d\n", k, v)
	}
}

func put(tdata *uint64, tblen *uint64, tbitlen *uint64) {
	var mask uint64
	for *tblen >= 8 {
		newbyte := byte((*tdata >> (*tblen - 8))&0xFF)
		bits = append(bits, newbyte)
		//fmt.Printf("put: tblen=%d, tbitlen=%d, tdata=0x%08x, newbyte=0x%02x\n", *tblen, *tbitlen, *tdata, newbyte)
		mask = 0xFF
		mask = ^(mask << (*tblen - 8))
		*tdata &= mask
		*tblen -= 8
		*tbitlen += 8
		//fmt.Printf("put: tblen=%d, tbitlen=%d, tdata=0x%08x, mask=0x%08x\n", *tblen, *tbitlen, *tdata, mask)
	}
}

func fill_random(nvalues int) (tbitlen uint64) {
var tdata	uint64	// buffer where bits are stored, drained to less than 8 bits, after data is added
var tblen	uint64	// used to keep track of how many bit used in the above buffer
var sav		uint64
var tmp		uint64

	bits = nil
	data = make([]info, nvalues)
	for k := range data {
		v := &data[k]
		for {
			v.blen = uint(rbetween(1, 32))
			if true { // v.blen % 4 == 0
				break
			}
		}
		v.value = uint32(rbetween(0, int(math.Exp2(float64(v.blen))-1)))

		//fmt.Printf("value=0x%x, blen=%d\n", v.value, v.blen)
		tdata <<= v.blen
		tdata |= uint64(v.value)
		tblen += uint64(v.blen)
		put(&tdata, &tblen, &tbitlen)
		dist[v.blen]++
		//fmt.Printf("k=%d, blen=%d, value=0x%x\n", k, v.blen, v.value)
	}
	sav = tblen
	//fmt.Printf("tblen=%d, ", tblen)
	for {
		if tblen == 0 || tblen >= 8 {
			break
		}
		tdata <<= 1 // fill with zeros
		tblen++
	}
	//fmt.Printf("%d bits filled\n", tblen - sav)
	sav++
	put(&tdata, &tblen, &tmp) // don't include fill in bitlen
	return
}

func fill_FF(nbytes int) (tbits int) {
	bits = make([]byte, nbytes)
	//fmt.Printf("nbytes=%d, len(bits)=%d\n", nbytes, len(bits))
	for i := range bits {
		bits[i] = 0xFF
	}
	return nbytes*8
}


func TestFF(t *testing.T) {
	var mbits = maxbitlen

	defer func() {
		if p := recover(); p != nil {
			if p == "EOF" {
				return
			}
			fmt.Printf("unknown error: %v", p)
			return
		}
	}()

	fmt.Printf("TestFF\n")
	bits = nil
	tbits := fill_FF(nBytes)
	bs, _ := NewFromMemory(bits)

	for nbits := tbits; nbits > 0; {
		if (nbits < mbits) {
			//fmt.Printf("Test_FF: mbits=%d, setting mbits to %d\n", mbits, nbits)
			mbits = nbits
		}
		blen := uint(rbetween(0, mbits))
		pvalue := bs.Peekbits(blen)
		value := bs.Getbits(blen)
		comp := uint32((1<<uint32(blen))-1)
		//fmt.Printf("Test_FF: nbits=%d, blen=%d, value=%d, value=0x%x\n", nbits, blen, value, value)

		if pvalue != comp {
			fmt.Printf("Test_FF: bs.Peekbits ERROR: blen=%d, value=0x%x, comp=0x%x\n", blen, pvalue, comp)
		}
		if value != comp {
			fmt.Printf("Test_FF: bs.Getbits  ERROR: blen=%d, value=0x%x, comp=0x%x\n", blen, value, comp)
		}

		nbits -= int(blen)
		//fmt.Printf("Test_FF: nbits=%d, blen=%d\n", nbits, blen)
	}
	fmt.Printf("Test_FF: tested tbits=%d\n", tbits)
}

func TestRandom(t *testing.T) {
	fmt.Printf("TestRandom\n")
	tbits := fill_random(nEntries)
	//dumpDist()
	bs, _ := NewFromMemory(bits)
	for k, v := range data {
		// fmt.Printf("value=0x%x, blen=%d\n", v.value, v.blen)
		pvalue := bs.Peekbits(v.blen)
		value := bs.Getbits(v.blen)
		if (pvalue != v.value) {
			fmt.Printf("TestRandom: bs.Peekbits ERROR: k=%d, blen=%d, v.value=0x%x, value=0x%x \n", k, v.blen, v.value, value)
		}
		if (value != v.value) {
			fmt.Printf("TestRandom: bs.Getbits  ERROR: k=%d, blen=%d, v.value=0x%x, value=0x%x \n", k, v.blen, v.value, value)
		}
	}
	fmt.Printf("TestRandom: tested %d entries with tbits=%d\n", len(data), tbits)
}

