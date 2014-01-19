package bitstream

import "fmt"
//import "flag"
import "os"
import "io"

type bitstream struct {
	file			*os.File	// file pointer of stream
	r				io.Reader	// interface to read bytes from
	buf				[]byte		// read buffer
	bp				[]byte
	bufn			uint32		// next buf, next 32 bits, this is always 32 bits
	nbits			uint			// next bits, either 0, 8, 16, 24, or 32
	bufc			uint32		// current buf, the actual bits, comes from bufn
	bits			uint			// the number of bits remaining in strm_bufc
	tbits			uint64		// the total number of bits read so far
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

func (bs *bitstream) Open(path string, perm os.FileMode) error {
	file, err := os.OpenFile(path, os.O_RDONLY, perm)
	if err != nil {
		fmt.Printf("bitstream.Open: can't open name=%q, perm=%o\n", path, perm)
		return err
	}
	bs.file = file
	return nil
}

func (bs *bitstream) Close() {
	fmt.Printf("bitstream.Close\n")
}

func Init() *bitstream {
var		bs bitstream

	bs.buf = make([]byte, 4)
	bs.bp = bs.buf[0:0]
	fmt.Printf("bitstream.Init: len(bs.buf)=%d\n", len(bs.bp))
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

func dump(buf []byte, max int) {
	for k := range buf {
		if k >= max {
			break
		}
		fmt.Printf("buf[%d]=0x%02x, ", k, buf[k])
	}
	fmt.Printf("\n")
}

func NewFromFile(path string) (*bitstream, error) {
	fmt.Printf("bitstream.NewFromFile: New path=%q\n", path)
	bs := Init()
	if err := bs.Open(path, 0666); err != nil {
		fmt.Printf("bitstream.NewFromFile: New Error\n")
		return nil, err
	}
	fmt.Printf("bitstream.NewFromFile: New Ok\n")
	bs.r = bs.file
	fmt.Printf("bitstream.NewFromFile: bs.r=%v, bs.r=%p\n", bs.r, bs.r)
	bs.readbits()
	return bs, nil
}

type Memory struct {
	buf		[]byte
}

// Read reads up to len(b) bytes from memory.
// It returns the number of bytes read.
// EOF is signaled by a zero count with err set to io.EOF.
func (m *Memory) Read(b []byte) (n int, err error) {
	if m == nil {
 		return 0, os.ErrInvalid
	}
//	fmt.Printf("Memory.Read: len(b)=%d, len(m.buf)=%d, cap(m.buf)=%d\n", len(b), len(m.buf), cap(m.buf))
    if len(m.buf) == 0 {
		fmt.Printf("Memory.Read: EOF\n")
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
//	fmt.Printf("Memory.Read: A len(m.buf)=%d, cap(m.buf)=%d\n", len(m.buf), cap(m.buf))
    return n, nil
}


func NewFromMemory(b []byte) (*bitstream, error) {
	fmt.Printf("bitstream.NewFromMemory\n")
	bs := Init()
	m := Memory{b}
	fmt.Printf("bitstream.New: New Ok\n")
	bs.r = &m
	fmt.Printf("bitstream.New: bs.r=%v, bs.r=%p\n", bs.r, bs.r)
	bs.readbits()
	return bs, nil
}

func (bs *bitstream) readbits() error {
	//fmt.Printf("bitstream.readbits: Read()\n")
	if bs.eof {
		return io.EOF
	}
	bs.bufn = 0
	bs.nbits = 0
	//fmt.Printf("bitstream.readbits: Read() 2\n")

	for i := 0; i < 4; i++ {
		if len(bs.bp) == 0 {
			bs.bp = bs.buf[0:cap(bs.buf)]
			//fmt.Printf("bitstream.readbits: Read() len(bs.buf)=%d\n", len(bs.bp))
			l, err := bs.r.Read(bs.bp[:])
			//fmt.Printf("bitstream.readbits: Read() l=%d, err=%v\n", l, err)
			if l < 0 || err != nil {
				bs.eof = true
				if (bs.nbits > 0) {
					return nil
				}
			}
	//		dump(bs.bp, 10)
		}

//		printf("b=0x%02x\n", tmp);
		bs.bufn <<= 8
		bs.bufn |= uint32(bs.bp[0])
		bs.bp = bs.bp[1:]
		//fmt.Printf("bitstream.readbits: len(bs.buf)=%d, cap(bs.buf)=%d\n", len(bs.bp), cap(bs.bp))
		bs.nbits += 8
	}
//	printf("sp->strm_nbits=%ld, strm_bufn=0x%08lx\n", sp->strm_nbits, sp->strm_bufn);
	return nil
}


func (bs *bitstream) getbits2(bits uint) (uint32, error) {
var rbits uint = bits
// var tmp uint
var ret uint32 = 0

	fmt.Printf("bitstream.getbits: bits=%d\n", bits);
	//fmt.Printf("bitstream.getbits: bs=%#v\n", bs);
	if bs.eof == true  && bs.bits == 0 && bs.nbits == 0 {
		return 0, io.EOF
	}
	
	if bits <= bs.bits {
		ret = ((bs.bufc>>(bs.bits - bits))&((1<<bits)-1))
		//fmt.Printf("bitstream.getbits2 1ret=0x%x\n", ret)
	} else {
		if bs.bits > 0 {
			rbits = bits - bs.bits
			ret = (bs.bufc&((1<<bs.bits)-1)) << rbits
			//fmt.Printf("bitstream.getbits2 2ret=0x%x\n", ret)
		}
		bs.bufc = bs.bufn
		bs.bits = bs.nbits
		if err := bs.readbits(); err != nil {
			bs.nbits = 0
			return 0, err
		}
		if (rbits > bs.bits) {
			// tmp = rbits - bs.bits;
			ret |= (bs.bufc&((1<<bs.bits)-1));  // we could zero fill, but we don't
			//fmt.Printf("bitstream.getbits2 3ret=0x%x\n", ret)
		} else {
			ret |= ((bs.bufc>>(bs.bits - rbits))&((1<<rbits)-1))
			//fmt.Printf("bitstream.getbits2 4ret=0x%x\n", ret)
		}
	}
	bs.bits -= bits
	bs.tbits += uint64(bits)
	return ret, nil
}


// same as above but just peek, don't update any counters
func (bs *bitstream) peekbits2(bits uint) (uint32, error) {
var rbits uint = bits
var ret uint32 = 0

	if bs.eof == true  && bs.bits == 0 && bs.nbits == 0 {
		return 0, io.EOF
	}
	
	if bits <= bs.bits {
		return (bs.bufc>>(bs.bits - bits))&((1<<bits)-1), nil
	} else {
		if bs.bits > 0 {
			rbits = bits - bs.bits
			ret = (bs.bufc&((1<<bs.bits)-1)) << rbits
			rbits = bits - bs.bits
		}
		return ret | ((bs.bufn>>(bs.bits - rbits))&((1<<rbits)-1)), nil
	}
	return 0, io.EOF // not really an EOF need generic error
}

func (bs *bitstream) peekbits(bits uint) uint32 {
	r, _ := bs.peekbits2(bits)
	return r
}

func (bs *bitstream) getbits(bits uint) uint32 {
	r, _ := bs.getbits2(bits)
	return r
}

func (bs *bitstream) get_byte_aligned() error {
	for (bs.tbits&0x7) != 0 {
		bs.rub()
	}
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

func (bs *bitstream) skipbits(bits uint) error {
	for bits > 32 {
		bs.rul()
		bits -= 32
	}
	return nil
}

// read unsigned long
func (bs *bitstream) rul() uint32 {
	tmp, _ := bs.getbits2(32)
	return tmp
}

	
// read unsigned short
func (bs *bitstream) rus() uint16 {
	ret, _ := bs.getbits2(16)
	ret = ret&0xFFFF
	return uint16(ret)
}


// read unsigned char
func (bs *bitstream) ruc() byte {
	ret, err := bs.getbits2(8)
	if (err != nil) {
		fmt.Printf("err=%v\n", err)
		panic("EOF")
	}
	ret = ret&0xFF
	fmt.Printf("bitstream.ruc ret=0x%02x\n", ret)
	return byte(ret)
}


// read bit or bool
func (bs *bitstream) rub() bool {

	ret, _ := bs.getbits2(1)
	ret = ret&0x1
	if ret == 1 {
		return true
	} else {
		return false
	}
}

// read unsigned long sub
func (bs *bitstream) ruls(bits uint) uint32 {
	ret, _ := bs.getbits2(bits)
	return ret
}

	
// read unsigned short sub
func (bs *bitstream) russ(bits uint) uint16 {

/*
	if (bits > 16)
		iexit("russ");
*/
//	printf("russ: 0x%lx\n", ul&0xFFFF);
	ret, _ := bs.getbits2(bits)
	return uint16(ret)
}

// read unsigned short sub
func (bs *bitstream) rss(bits uint) int16 {

/*
	if (bits > 16)
		iexit("russ");
*/
//	printf("russ: 0x%lx\n", ul&0xFFFF);
	ret, _ := bs.getbits2(bits)
	return int16(ret)
}


// read unsigned char sub
func (bs *bitstream) rucs(bits uint) byte {
	ret, _ := bs.getbits2(bits)
	return byte(ret&0xFF)
}

// read signed char sub
func (bs *bitstream) rcs(bits uint) int8 {
	ret, _ := bs.getbits2(bits)
	return int8(ret)
}