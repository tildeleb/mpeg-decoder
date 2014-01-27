package bitstream

import "fmt"
//import "flag"
import "os"
import "io"

// this package implements reading from a steam of bytes 1-32 bits at a time
// various routines are provided to read or peek 1-32 bits and to read a variable
// number of bits into an int8, int16, int32, or int64 or their unsigned counterparts
// basic stragey is to maintain 2 x 32 bit buffers bufc (current) and bufn (next)
// bufc is a shift register, bit come off the MSB side because MPEG is big-endian based
// get bits from bufc, if not enought bits are available, copy bufn to bufc and repeat
// we typically fetch 4 bytes, one at a time to be divorced from endianess
// still unoptimized speed appears to be about 250 Mbit/sec on my 2.5 GHz i7
type Bitstream struct {
	file			*os.File	// file pointer of stream
	r				io.Reader	// interface to read bytes from
	buf				[]byte		// read buffer
	bp				[]byte		// slice into read buffer
	bufn			uint32		// next buf, next 32 bits usually, number of bits defined by nbits
	nbits			uint		// next bits, either 0, 8, 16, 24, or 32
	bufc			uint32		// current buf, the next bits come from here. If not enough copy bufn to here (bufc)
	bits			uint		// the number of bits remaining in bufc can be 1-32, if 0 we get more
	rbits			uint64		// the total number of bits read so far
	tbits			uint64		// total bits available
	eof				bool		// have reached eof, no more bits to be had
}


/*
bool
open_file(char* np, FILE** fpp, char* extp, char* modep)
{
char		name[256];
char*		cp;

	strcpy(name, np);
	if (extp != 0)
		if ( (cp = strchr(np, '.')) != 0)
			strcpy(cp, extp);
	if ( (*fpp = fopen(name, modep)) == NULL)
		return(true);
	return(false);
}
*/

func (bs *Bitstream) Open(path string, perm os.FileMode) error {
	file, err := os.OpenFile(path, os.O_RDONLY, perm)
	if err != nil {
		fmt.Printf("bitstream.Open: can't open name=%q, perm=%o\n", path, perm)
		return err
	}
	bs.file = file
	return nil
}

func (bs *Bitstream) Close() {
	fmt.Printf("bitstream.Close\n")
}

func Init() *Bitstream {
var		bs Bitstream

	bs.buf = make([]byte, 4)
	bs.bp = bs.buf[0:0]
	//fmt.Printf("bitstream.Init: len(bs.buf)=%d\n", len(bs.bp))
	bs.bufn = 0
	bs.bufc = 0
	bs.bits = 0
	bs.nbits = 0
	bs.eof = false
	return &bs
}

/*
func NewReader(r io.Reader) *bitstream {
	bs := Init()
	bs.r = r
}
*/

func Dump(buf []byte, max ...int) {
	for k := range buf {
		if len(max) > 0 && k >= max[0] {
			break
		}
		// fmt.Printf("buf[%d]=0x%02x, ", k, buf[k])
		fmt.Printf("0x%02x, ", buf[k])
	}
	fmt.Printf("\n")
}

func NewFromFile(path string) (*Bitstream, error) {
	fmt.Printf("bitstream.NewFromFile: New path=%q\n", path)
	bs := Init()
	if err := bs.Open(path, 0666); err != nil {
		fmt.Printf("bitstream.NewFromFile: New Error\n")
		return nil, err
	}
	fmt.Printf("bitstream.NewFromFile: New Ok\n")
	bs.r = bs.file
	//fmt.Printf("bitstream.NewFromFile: bs.r=%v, bs.r=%p\n", bs.r, bs.r)
	bs.readbits()
	return bs, nil
}

type Memory struct {
	buf		[]byte
}

func NewFromMemory(b []byte) (*Bitstream, error) {
	//fmt.Printf("bitstream.NewFromMemory\n")
	bs := Init()
	bs.tbits = uint64(len(b) * 8)
	m := Memory{b}
	//fmt.Printf("bitstream.New: New Ok\n")
	bs.r = &m
	//fmt.Printf("bitstream.New: bs.r=%v, bs.r=%p\n", bs.r, bs.r)
	bs.readbits()
	return bs, nil
}

// Read reads up to len(b) bytes from memory.
// It returns the number of bytes read.
// EOF is signaled by a zero count with err set to io.EOF.
func (m *Memory) Read(b []byte) (n int, err error) {
	if m == nil {
 		return 0, os.ErrInvalid
	}
	//fmt.Printf("Memory.Read: len(b)=%d, len(m.buf)=%d, cap(m.buf)=%d\n", len(b), len(m.buf), cap(m.buf))
    if len(m.buf) == 0 {
		//fmt.Printf("Memory.Read: EOF\n")
    	return 0, io.EOF
    }
    if len(b) > len(m.buf) {
    	n = len(m.buf)
    } else {
    	n = len(b)
    }
//	fmt.Printf("Memory.Read: B len(m.buf)=%d, cap(m.buf)=%d\n", len(m.buf), cap(m.buf))
	for i := 0; i < n; i++ {
		b[i] = m.buf[i]
	}
//    b = m.buf[0:n]
    m.buf = m.buf[n:]
	//fmt.Printf("Memory.Read: A len(m.buf)=%d, cap(m.buf)=%d\n", len(m.buf), cap(m.buf))
    return n, nil
}

func (bs *Bitstream) readbits() error {
	//fmt.Printf("bitstream.readbits: Read()\n")
	if bs.eof {
		return io.EOF
	}
	bs.bufn = 0
	bs.nbits = 0
	//fmt.Printf("bitstream.readbits: Read() 2\n")

	cnt := 4
	for {
		if len(bs.bp) == 0 {
			bs.bp = bs.buf[0:cap(bs.buf)]
			//fmt.Printf("bitstream.readbits: Read() bs.nbits=%d, len(bs.buf)=%d\n", bs.nbits, len(bs.bp))
			l, err := bs.r.Read(bs.bp[:])
			//fmt.Printf("bitstream.readbits: Read() l=%d, err=%v\n", l, err)
			if l <= 0 || err != nil {
				bs.eof = true
				// if (bs.nbits > 0) {
					return nil
				// }
			}
			bs.bp = bs.bp[:l]
			if (l < cnt) {
				cnt = l
			}
	//		dump(bs.bp, 10)
		}

		//fmt.Printf("bitstream.readbits: bs.bp[0]=0x%02x\n", bs.bp[0])
		bs.bufn <<= 8
		bs.bufn |= uint32(bs.bp[0])
		bs.bp = bs.bp[1:]
		bs.nbits += 8
		cnt--
		//fmt.Printf("bitstream.readbits: nbits=%d, len(bs.buf)=%d, cap(bs.buf)=%d\n", bs.nbits, len(bs.bp), cap(bs.bp))
		if cnt == 0 {
			break
		}
	}
//	printf("sp->strm_nbits=%ld, strm_bufn=0x%08lx\n", sp->strm_nbits, sp->strm_bufn);
	return nil
}


func (bs *Bitstream) getbits2(bits uint) (uint32, error) {
var rbits uint = bits
// var tmp uint
var ret uint32
var fbits uint

	//fmt.Printf("bitstream.getbits2: rbits=%d/tbits=%d, bits=%d, nbits=%d, eof=%v, rbits=%d\n", bs.rbits, bs.tbits, bs.bits, bs.nbits, bs.eof, bits);
	//fmt.Printf("bitstream.getbits2: bs.bufc=0x%x, bs.bits=%d, bs.bufn=0x%x bs.nbits=%d\n", bs.bufc, bs.bits, bs.bufn, bs.nbits)
	//fmt.Printf("bitstream.getbits: bs=%#v\n", bs);
	if bs.eof == true  && bs.bits == 0 && bs.nbits == 0 {
		os.Exit(0)
		panic("EOF")
		// return 0, io.EOF
	}

	// we fill bits if a request goes beyond the EOF mainly to make automaked testing work better
	// none of the primitives return the number of bits actaully read, which would be a pain and
	// this isn't required for MPEG streams
	if (bits > bs.bits + bs.nbits) {
		fbits = bits - (bs.bits + bs.nbits)
		bits = bs.bits + bs.nbits
	}
	
	if bits <= bs.bits {
		ret = ((bs.bufc>>(bs.bits - bits))&((1<<bits)-1))
		bs.bits -= bits
		//fmt.Printf("bitstream.getbits2 1ret=0x%x\n", ret)
	} else {
		if bs.bits > 0 {
			rbits = bits - bs.bits
			ret = (bs.bufc&((1<<bs.bits)-1)) << rbits
			//fmt.Printf("bitstream.getbits2 2ret=0x%x\n", ret)
		}
		bs.bufc = bs.bufn
		bs.bits = bs.nbits
		bs.nbits = 0
		if err := bs.readbits(); err != nil {
			bs.nbits = 0
			//return 0, err
		}
		if (rbits > bs.bits) {
			// tmp = rbits - bs.bits;
			ret |= (bs.bufc&((1<<bs.bits)-1));  // we could zero fill, but we don't
			//fmt.Printf("bitstream.getbits2 3ret=0x%x\n", ret)
		} else {
			ret |= ((bs.bufc>>(bs.bits - rbits))&((1<<rbits)-1))
			bs.bits -= rbits
			//fmt.Printf("bitstream.getbits2 4ret=0x%x\n", ret)
		}
	}
	fbits++
/*
	if fbits > 0 {
		ret <<= fbits
	}
*/
	bs.rbits += uint64(bits)
	//fmt.Printf("bitstream.getbits2: rbits=%d/tbits=%d, bits=%d, nbits=%d, eof=%v, rbits=%d, ret=0x%x\n", bs.rbits, bs.tbits, bs.bits, bs.nbits, bs.eof, bits, ret);
	return ret, nil
}


// same as above but just peek, don't update any counters
func (bs *Bitstream) peekbits2(bits uint) (uint32, error) {
var rbits uint = bits
var ret uint32 = 0

	//fmt.Printf("bitstream.peekbits2: rbits=%d/tbits=%d, bits=%d, nbits=%d, eof=%v, rbits=%d\n", bs.rbits, bs.tbits, bs.bits, bs.nbits, bs.eof, bits);
	//fmt.Printf("bitstream.peekbits2: bs.bufc=0x%x, bs.bits=%d, bs.bufn=0x%x bs.nbits=%d\n", bs.bufc, bs.bits, bs.bufn, bs.nbits)

	//fmt.Printf("bitstream.peekbits2: bits=%d, nbits=%d, eof=%v, rbits=%d\n", bs.bits, bs.nbits, bs.eof, bits);
	if bs.eof == true  && bs.bits == 0 && bs.nbits == 0 {
		return 0, io.EOF
	}
	
	if bits <= bs.bits {
		ret = (bs.bufc>>(bs.bits - bits))&((1<<bits)-1)
	} else {
		if bs.bits > 0 {
			rbits = bits - bs.bits
			ret = (bs.bufc&((1<<bs.bits)-1)) << rbits
		}

		ret |= ((bs.bufn>>(bs.nbits - rbits))&((1<<rbits)-1))
	}
	//fmt.Printf("bitstream.peekbits2: bits=%d, nbits=%d, eof=%v, rbits=%d, ret=0x%x\n", bs.bits, bs.nbits, bs.eof, bits, ret);
	//fmt.Printf("bitstream.peekbits2: rbits=%d/tbits=%d, bits=%d, nbits=%d, eof=%v, rbits=%d, ret=0x%x\n", bs.rbits, bs.tbits, bs.bits, bs.nbits, bs.eof, bits, ret);
	return ret, nil
}

func (bs *Bitstream) Peekbits(bits uint) uint32 {
	r, _ := bs.peekbits2(bits)
	//fmt.Printf("Peekbits: bits=%d, ret=0x%x\n", bits, r)
	return r
}

func (bs *Bitstream) Getbits(bits uint) uint32 {
	r, _ := bs.getbits2(bits)
	//fmt.Printf("Getbits: bits=%d, ret=0x%x\n", bits, r)
	return r
}

func (bs *Bitstream) GetByteAligned() error {
	for (bs.rbits&0x7) != 0 {
		bs.Rub()
	}
	fmt.Printf("GetByteAligned: tbits=0x%x\n", bs.rbits)
	return nil
}

/*
				sp->strm_buf =	((tmp>>24)&0xFF) | (((tmp>>16)&0xFF)<<8) |
								(((tmp>>8)&0xFF)<<16) | ((tmp&0xFF)<<24);

				sp->strm_buf =	(((tmp>>16)&0xFF)<<8) |
								(((tmp>>8)&0xFF)<<16) | ((tmp&0xFF)<<24);

				sp->strm_buf = (((tmp>>8)&0xFF)<<16) | ((tmp&0xFF)<<24);

				sp->strm_buf = ((tmp&0xFF)<<24);
*/

func (bs *Bitstream) Skipbits(bits uint) error {
	for bits > 32 {
		bs.Rul()
		bits -= 32
	}
	for bits > 0 {
		bs.Rub()
		bits -= 1
	}
	return nil
}


// read long
func (bs *Bitstream) Rl() int32 {
	tmp, _ := bs.getbits2(32)
	return int32(tmp)
}

	
// read short
func (bs *Bitstream) Rs() int16 {
	ret, _ := bs.getbits2(16)
	ret = ret&0xFFFF
	return int16(ret)
}


// read char
func (bs *Bitstream) Rc() byte {
	ret, err := bs.getbits2(8)
	if (err != nil) {
		fmt.Printf("err=%v\n", err)
		// panic("EOF")
	}
	ret = ret&0xFF
	fmt.Printf("bitstream.ruc ret=0x%02x\n", ret)
	return byte(ret)
}

// read unsigned long
func (bs *Bitstream) Rul() uint32 {
	tmp, _ := bs.getbits2(32)
	return tmp
}

	
// read unsigned short
func (bs *Bitstream) Rus() uint16 {
	ret, _ := bs.getbits2(16)
	ret = ret&0xFFFF
	return uint16(ret)
}


// read unsigned char
func (bs *Bitstream) Ruc() byte {
	ret, err := bs.getbits2(8)
	if (err != nil) {
		fmt.Printf("err=%v\n", err)
		// panic("EOF")
	}
	ret = ret&0xFF
	fmt.Printf("bitstream.ruc ret=0x%02x\n", ret)
	return byte(ret)
}


// read bit or bool
func (bs *Bitstream) Rub() bool {

	ret, _ := bs.getbits2(1)
	ret = ret&0x1
	if ret == 1 {
		return true
	} else {
		return false
	}
}

// read unsigned long sub
func (bs *Bitstream) Ruls(bits uint) uint32 {
	ret, _ := bs.getbits2(bits)
	return ret
}

	
// read unsigned short sub
func (bs *Bitstream) Russ(bits uint) uint16 {

/*
	if (bits > 16)
		iexit("russ");
*/
//	printf("russ: 0x%lx\n", ul&0xFFFF);
	ret, _ := bs.getbits2(bits)
	return uint16(ret)
}

// read unsigned short sub
func (bs *Bitstream) Rss(bits uint) int16 {

/*
	if (bits > 16)
		iexit("russ");
*/
//	printf("russ: 0x%lx\n", ul&0xFFFF);
	ret, _ := bs.getbits2(bits)
	return int16(ret)
}


// read unsigned char sub
func (bs *Bitstream) Rucs(bits uint) byte {
	ret, _ := bs.getbits2(bits)
	return byte(ret&0xFF)
}

// read signed char sub
func (bs *Bitstream) Rcs(bits uint) int8 {
	ret, _ := bs.getbits2(bits)
	return int8(ret)
}