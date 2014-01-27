package iso11172	// iso111722 rename on import

// Copyright © 2003 and 2014 Lawrence E. Bakst. All rights reserved.
// THIS SOURCE CODE IS THE PROPRIETARY INTELLECTUAL PROPERTY AND CONFIDENTIAL
// INFORMATION OF LAWRENCE E. BAKST AND IS PROTECTED UNDER U.S. AND
// INTERNATIONAL LAW. ANY USE OF THIS SOURCE CODE WITHOUT THE PRIOR WRITTEN
// AUTHORIZATION OF LAWRENCE E. BAKST IS STRICTLY PROHIBITED.

// ISO-11172-2 data structures in Go
// transliterated from C versions written in 2003 so please excuse the C style naming 
// data from an mpeg-1 stream is unpacked into these data structures
// if used for an encoder customer encoder mashalling code would need to be written
// this could be extended to mpeg-2

import . "leb/mpeg-decoder/bitstream"

import "fmt"
//import "flag"
//import "os"
//import "io"

type Mpeg1Parse interface{
	ReadSeqenceHeader() *SequenceHeader
	ReadGroupHeader() *GroupHeader
	ReadPictureHeader() *PictureHeader
	ReadSliceHeader(uint32) *SliceHeader
	ReadMBAI() uint32
	ReadMBType(PictureType) (uint32, uint32, uint32, uint32, uint32)
	ReadMBMVM() int16 // read macro block motion vector m
	ReadYCbCr() (uint32, uint32, uint32, uint32, uint32, uint32)
	ReadDCDC()
	ReadMBDCTDCY()
	ReadMBDCTDCC()
	ReadMacroBlock(i int)
	ReadMacroBlocks()
	ReadMPEG1Steam()

	SetMBT(uint32, uint32, uint32, uint32, uint32)
	SetYCbCr(uint32, uint32, uint32, uint32, uint32, uint32)
}

type Block [8][8]byte

type PictureType int

const (
	pt_forbid PictureType = 0
	pt_ipict = 1
	pt_ppict = 2
	pt_bpict = 3
	pt_dpict = 4
)

var pt_str []string = []string{"0 pict", "I pict", "P pict", "B pict", "D pict"}

var b3_to_PT []PictureType = []PictureType{pt_forbid, pt_ipict, pt_ppict, pt_bpict, pt_dpict, pt_forbid, pt_forbid, pt_forbid}

/* Video codes */
const PICTURE_START_CODE		= 0x100
const LOWEST_SLICE_CODE			= 0x101
const HIGHEST_SLICE_CODE		= 0x1AF
const SLICE_MASK				= 0xFF
const USER_DATA_START_CODE		= 0x1B2
const SEQ_HEADER_CODE			= 0x1B3
const SEQ_ERROR_CODE			= 0x1B4
const EXTENSION_START_CODE		= 0x1B5
const RESERVED_CODE				= 0x1B6
const SEQ_END_CODE				= 0x1B7
const GROUP_START_CODE			= 0x1B8



/* System codes */
const ISO_11172_END_CODE		= 0x1B9		
const PACK_START_CODE			= 0x1BA
const SYSTEM_HEADER_START_CODE	= 0x1BB
const PACKET_START_CODE_MASK	= 0xffffff00
const PACKET_START_CODE_PREFIX	= 0x00000100

const EOB						= 0x2


type SequenceHeader struct {
	sh_code							uint32
	sh_hor_size						uint16
	sh_ver_size						uint16
	sh_pel_aspect_ratio				uint16
	sh_picture_rate					uint16
	sh_bit_rate						uint32
	sh_marker_bit					bool
	sh_vbv_buffer_size				uint32
	sh_const_params_flag			bool
	sh_load_intra_quant_matrix		bool
	sh_load_non_intra_quant_matrix	bool
	sh_iqmp							*Block
	sh_niqmp						*Block
	sh_extp							[]byte
	sh_udp							[]byte
	sh_ghs							[]*GroupHeader
	sh_phs							[]*PictureHeader // some streams don't have GOPs
}

type GroupHeader struct {
	gh_code				uint32
//	gh_tc				uint32
	gh_tc_hr			byte
	gh_tc_min			byte
	gh_tc_sec			byte
	gh_tc_pic			byte
	gh_drop_frame_flag	bool
	gh_marker_bit		bool // always 1
	gh_closed_gap		bool
	gh_broken_link		bool
	gh_extp				[]byte
	gh_udp				[]byte
	gh_phs				[]*PictureHeader
}

type PictureHeader struct {
	ph_code					uint32
	ph_temporal_ref			uint16
	ph_picture_type			PictureType
	ph_vbv_delay			uint16
	ph_full_pell_forw_vec	bool
	ph_forw_code			int8
	ph_full_pell_back_vec	bool
	ph_back_code			int8
	ph_eip_count			int
	ph_eip					[]byte
	ph_ext					[]byte
	ph_ud					[]byte
	ph_shs					[]*SliceHeader
}

type SliceHeader struct {
	sl_code					uint32
	sl_quant_scale			byte
	sl_eip_count			int
	sl_eip					[]byte
	sl_mbh					[]*MacroBlockHeader
}

type MacroBlockHeader struct {
	mbt_pt					PictureType			
	mbt_ai					uint32
	mbt_qs					uint16	// quantizer scale
	
	mbt_in					bool
	mbt_pa					bool
	mbt_mb					bool
	mbt_mf					bool
	mbt_qf					bool
	
	mbt_mfhp				uint16
	mbt_mfhr				uint16
	mbt_mfvp				uint16
	mbt_mfvr				uint16
	mbt_mbhp				uint16
	mbt_mbhr				uint16
	mbt_mbvp				uint16
	mbt_mbvr				uint16
	
	mbt_blockv				[6]bool
	mbt_blockx				[6]Block
}

type MpegState struct {
	*Bitstream
	sh			[]*SequenceHeader
}

var next_code	uint32

func (ms *MpegState) ReadSeqenceHeader() *SequenceHeader {
var sh			SequenceHeader

//	b |= rul(sp, &shp->sh_code);
	fmt.Printf("bitstream.ReadSeqenceHeader: start\n")
	sh.sh_code = SEQ_HEADER_CODE
	sh.sh_hor_size = ms.Russ(12)
	sh.sh_ver_size = ms.Russ(12)
	sh.sh_pel_aspect_ratio = ms.Russ(4)
	sh.sh_picture_rate = ms.Russ(4)
	sh.sh_bit_rate = ms.Ruls(18)
	sh.sh_marker_bit = ms.Rub()
	sh.sh_vbv_buffer_size = ms.Ruls(10)
	sh.sh_const_params_flag = ms.Rub()
	sh.sh_load_intra_quant_matrix = ms.Rub()
	if sh.sh_load_intra_quant_matrix {
		fmt.Printf("bitstream.ReadSeqenceHeader: 1skip=%d\n", 8*64)
		ms.Skipbits(8*64)
	}
	sh.sh_load_non_intra_quant_matrix = ms.Rub()
	if sh.sh_load_non_intra_quant_matrix {
		fmt.Printf("bitstream.ReadSeqenceHeader: 2skip=%d\n", 8*64)
		ms.Skipbits(8*64)
	}
	sh.sh_iqmp = new(Block)
	sh.sh_niqmp = new(Block)
	sh.sh_extp = nil
	sh.sh_udp = nil
	fmt.Printf("bitstream.ReadSeqenceHeader: end\n")
	return &sh
}


func (ms *MpegState) ReadGroupHeader() *GroupHeader {
var gh			GroupHeader

	fmt.Printf("bitstream.ReadGroupHeader\n")
//	b |= rul(sp, &shp->sh_code);
	gh.gh_code = GROUP_START_CODE;
	gh.gh_drop_frame_flag = ms.Rub()
//	b |= ruls(sp, 25, &ghp->gh_tc);
	gh.gh_tc_hr = ms.Rucs(5)
	gh.gh_tc_min = ms.Rucs(6)
	gh.gh_marker_bit = ms.Rub()
	gh.gh_tc_sec = ms.Rucs(6)
	gh.gh_tc_pic = ms.Rucs(6)
	gh.gh_closed_gap = ms.Rub()
	gh.gh_broken_link = ms.Rub()
	gh.gh_extp = nil
	gh.gh_udp = nil
	return &gh
}


func (ms *MpegState) getExt() ([]byte, int) {
var abit		bool
var	cnt			int
var buf 		[100]byte
var bp			[]byte = buf[:]

	fmt.Printf("bitstream.getExt\n")
	for {
		abit = ms.Rub()
		if (abit) {
			bp[cnt] = ms.Ruc()
			cnt++
		} else {
			break
		}
		if (cnt >= len(buf)) {
			panic("getExt: out of space")
		}
	}
	bp = buf[:cnt]
	return bp, cnt
}


func (ms *MpegState) ReadPictureHeader() *PictureHeader {
//var abit		bool
//var	i, cnt		int
var uc			byte
//var buf 		[100]byte
//var bp			[]byte = buf[:]
var ph			PictureHeader


	fmt.Printf("bitstream.ReadPictureHeader\n")
//	b |= rul(sp, &shp->sh_code);
	ph.ph_code = PICTURE_START_CODE
	ph.ph_temporal_ref = ms.Russ(10)
	uc = ms.Rucs(3)
	ph.ph_picture_type = b3_to_PT[uc]
	ph.ph_vbv_delay = ms.Rus()
	if ph.ph_picture_type == pt_ppict || ph.ph_picture_type == pt_bpict {
		ph.ph_full_pell_forw_vec= ms.Rub()
		ph.ph_forw_code = ms.Rcs(3)
	}
	if ph.ph_picture_type == pt_bpict {
		ph.ph_full_pell_forw_vec = ms.Rub()
		ph.ph_forw_code = ms.Rcs(3)
	}
	
	ph.ph_eip, ph.ph_eip_count = ms.getExt()
	ph.ph_ext = nil;
	ph.ph_ud = nil;
	return &ph
}


func (ms *MpegState) ReadSliceHeader(code uint32) *SliceHeader {
//var abit bool
//var i, cnt int
//var buf 		[100]byte
//var bp			[]byte = buf[:]
var sl			SliceHeader
	
	fmt.Printf("bitstream.ReadSliceHeader\n")
	sl.sl_code = code;
	sl.sl_quant_scale = ms.Rucs(5)
	sl.sl_eip, sl.sl_eip_count = ms.getExt()
	return &sl
}


func (ms *MpegState) ReadMBAI() (ret uint32) {
var 			bits4a, bits4b, bits3, bits2, bits1 uint32

	// first try the three most common cases 0b1, 0b011, 0b010
	bits1 = ms.Peekbits(1)
	if bits1 == 1 {
		_ = ms.Getbits(1)
		ret = 1
	} else {
		bits3 = ms.Peekbits(3)
		if bits3 == 0x3 {
			_ = ms.Getbits(3)
			ret= 2
		} else {
			bits3 = ms.Peekbits(3)
			if bits3 == 0x2 {
				_ = ms.Getbits(3)
				ret = 3
			} else {
				// ok, it's not a common case, now we basically peel off 4 bits and then determine
				// how many more bits will need to get the whole code which is 4-11 bits
				bits4a = ms.Peekbits(4)
				switch bits4a {
				case 0:
					_ = ms.Getbits(4)
					bits4b = ms.Peekbits(4)
					switch bits4b {
					case 0:
						panic("read_mb_ai: zero stuffed")
						break;
					case 1:
						panic("read_mb_ai: bad pat")
						break;
					case 3:
						_ = ms.Getbits(4)
						bits3 = ms.Getbits(3)
						ret = 26 + (7 - bits3)
					case 4:
						_ = ms.Getbits(4)
						bits2 = ms.Peekbits(2)
						switch bits2 {
						case 0x2:
						case 0x3:
							_ = ms.Getbits(2)
							ret = 20 + (3 - bits2)
						case 0x00:
						case 0x01:
							_ = ms.Getbits(2)
							ret = 22 + (3 - bits2)
						default:
							panic("read_mb_ai: base 4");
						}
					case 5:
						_ = ms.Getbits(4)
						bits2 = ms.Getbits(2)
						ret = 16 + (3 - bits2)
					case 6, 7, 8, 9, 10, 11:
						_ = ms.Getbits(4)
						bits4b = ms.Getbits(4)
						ret = 10 + (11 - bits4b)
					default:
						panic("read_mb_ai: bad second 4 bits")
					}
				case 2:
				case 3:
					_ = ms.Getbits(4)
					ret = 4 + (3 - bits4a)
				case 1:
					_ = ms.Getbits(4)
					bits1 = ms.Getbits(1)
					ret = 6 + (1 - bits1)
				default:
					panic("read_mb_ai: bad first 4 bits");
				}
			}
		}
	}
	fmt.Printf("ReadMBAI: ret=%d\n", ret)
	return
}


func (ms *MpegState) SetMBT(mbh *MacroBlockHeader, in, pa, mb, mf, qf uint32) {
var ui32_to_bool = func(b uint32) bool {
	switch b {
	case 0:
		return false
	case 1:
		return true
	default:
		panic("ui32_to_bool: bad value")
	}
}

	fmt.Printf("iso.SetMBT in=%d, pa=%d, mb=%d, mf=%d, qf=%d\n", in, pa, mb, mf, qf)
	mbh.mbt_in = ui32_to_bool(in)
	mbh.mbt_pa = ui32_to_bool(pa)
	mbh.mbt_mb = ui32_to_bool(mb)
	mbh.mbt_mf = ui32_to_bool(mf)
	mbh.mbt_qf = ui32_to_bool(qf)
}


func (ms *MpegState) ReadMBType(pt PictureType) (in, pa, mb, mf, qf uint32) {
var			bits6 uint32
var			bits5 uint32
var			bits4 uint32
var			bits3 uint32
var			bits2 uint32
var			bits1 uint32
	
	bits1 = ms.Peekbits(1)
	if (bits1 == 1) {
		_ = ms.Getbits(1)
		switch (pt) {
		case pt_ipict:
			in, pa, mb, mf, qf = 1, 0, 0, 0, 0
		case pt_ppict:
			in, pa, mb, mf, qf = 0, 1, 0, 1, 0
		case pt_dpict:
			in, pa, mb, mf, qf = 1, 0, 0, 0, 0
		default:
			fmt.Printf("ReadMBType: mbt.mbt_pt=%d\n", pt);
			panic("ReadMBType: bad 1");
		}
		return
	}
	bits2 = ms.Peekbits(2)
	switch bits2 {
	case 1:
		_ = ms.Getbits(2)
		switch pt {
		case pt_ipict:
			in, pa, mb, mf, qf = 1, 0, 0, 0, 1
		case pt_ppict:
			in, pa, mb, mf, qf = 0, 1, 0, 0, 0
		default:
			panic("read_mb_type: bad 2");
		}
		return
	case 2:
		_ = ms.Getbits(2)
		if pt == pt_ppict {
			in, pa, mb, mf, qf = 0, 0, 1, 1, 0
		} else {
			panic("read_mb_type: bad 3")
		}
	case 3:
		_ = ms.Getbits(2)
		if pt == pt_ppict {
			in, pa, mb, mf, qf = 0, 1, 1, 1, 0
		} else {
			panic("read_mb_type: bad 4")
		}
	}
	bits3 = ms.Peekbits(3)
	switch(bits3) {
	case 1:
		_ = ms.Getbits(3)
		switch pt {
		case pt_ipict:
			in, pa, mb, mf, qf = 1, 0, 0, 0, 1
		case pt_ppict:
			in, pa, mb, mf, qf = 0, 1, 0, 0, 0
		default:
			panic("read_mb_type: bad 2")
		}
		return
	case 2:
		_ = ms.Getbits(3) // ???
		if pt == pt_ppict {
			in, pa, mb, mf, qf = 0, 0, 1, 1, 0
		} else {
			panic("read_mb_type: bad 3")
		}
	case 3:
		_ = ms.Getbits(3) // ???
		if pt == pt_ppict {
			in, pa, mb, mf, qf = 0, 1, 1, 1, 0
		} else {
			panic("read_mb_type: bad 4")
		}
	}
	bits4 = ms.Peekbits(4)
	switch bits4 {
	case 0:
		if pt == pt_ppict {	// only 1 5 bit code that starts with 0b0000
			bits5 = ms.Peekbits(5)
			if bits5 == 1 {
				bits5 = ms.Getbits(5)
				in, pa, mb, mf, qf = 0, 1, 0, 0, 1
			}
		}
		bits6 = ms.Getbits(6)
		switch(bits6) {
		case 0x01:
			if pt == pt_ppict {
				in, pa, mb, mf, qf = 1, 0, 0, 0, 1
			} else {
				if pt == pt_bpict {
					in, pa, mb, mf, qf = 1, 0, 0, 0, 1
				} else {
					panic("read_mb_type: bad 0b000001")
				}
			}
		case 0x02:
			if pt != pt_bpict {
				panic("read_mb_type: bad 0b000010")
				in, pa, mb, mf, qf = 0, 1, 1, 0, 1
			} else {
				panic("read_mb_type: 0b000010 expected pt_bpict")
			}
		case 0x03:
			if pt != pt_bpict {
				panic("read_mb_type: bad 0b000011")
				in, pa, mb, mf, qf = 0, 1, 0, 1, 1
			} else {
				panic("read_mb_type: 0b000011 expected pt_bpict")
			}
		default:
			panic("read_mb_type: bad 6 bit code")
		}
	case 0x01:	// all 5 bits codes
		bits5 = ms.Getbits(5)
		switch bits5 {
		case 0x2:
			if pt == pt_ppict {
				in, pa, mb, mf, qf = 0, 1, 0, 1, 1
			} else {
				if pt == pt_bpict {
					in, pa, mb, mf, qf = 0, 1, 1, 1, 1
				} else {
					panic("read_mb_type: bad 0b00010")
				}
			}
		case 0x03:
			if pt == pt_ppict {
			in, pa, mb, mf, qf = 1, 0, 0, 0, 0
			} else {
				if pt == pt_bpict {
					in, pa, mb, mf, qf = 1, 0, 0, 0, 0
				} else {
					panic("read_mb_type: bad 0b00011")
				}
			}
		default:
			panic("read_mb_type: bad 5 bit code")
		}
	case 0x02:
		bits4 = ms.Getbits(4)
		if pt != pt_ppict {
			panic("read_mb_type: code 0b0010 expected pt_ppict")
		}
		in, pa, mb, mf, qf = 0, 0, 0, 1, 0
	case 0x03:
		bits4 = ms.Getbits(4)
		if pt != pt_ppict {
			panic("read_mb_type: code 0b0010 expected pt_ppict")
		}
		in, pa, mb, mf, qf = 0, 1, 0, 1, 0
	default:
		panic("read_mb_type: bad 4 bit code")
	}
	return
}


func (ms *MpegState) ReadMBMVM() int16 {
var bits4a	uint32
var bits4b	uint32
var bits3	uint32
var bits2	uint32
var bits1	uint32
var ternary = func(c bool, a, b int16) int16 {
	if (c) {
		return a
	} else {
		return b
	}
}

	bits1 = ms.Peekbits(1)
	if bits1 == 1 {
		_ = ms.Getbits(1)
		return 0
	}
	bits3 = ms.Peekbits(3)
	if (bits3&0x6) == 0x010 {
		_ = ms.Getbits(3)
		return ternary((bits3&0x1) == 1, -1, 1)
	}
	bits4a = ms.Peekbits(4)
	switch(bits4a) {
	case 0x2:
	case 0x3:
		_ = ms.Getbits(4)
		return ternary((bits4a&0x01) == 1, -2, 2)
	case 0x1:
		_ = ms.Getbits(4)
		bits1 = ms.Getbits(1)
		return ternary((bits1&0x01) == 1, -3, 3)
	case 0x0:
		// ???
		_ = ms.Getbits(4)
		bits3 = ms.Peekbits(3)
		if (bits3&0x6) == 0x6 {
			_ = ms.Getbits(3)
			return ternary((bits3&0x01) == 1, -4, 4)
		}
		// guaranteed to have 4 bits now, this is the second set of 4 bits
		bits4b = ms.Peekbits(4)
		switch bits4b {
		case 0xA:
		case 0xB:
			return ternary((bits4b&0x01) == 1, -5, 5)
		case 0x8:
		case 0x9:
			return ternary((bits4b&0x01) == 1, -6, 6)
		case 0x6:
		case 0x7:
			return ternary((bits4b&0x01) == 1, -7, 7)
		case 0x5:
			bits2 = ms.Getbits(2)
			switch bits2 {
			case 2:
				return 8
			case 3:
				return -8
			case 0:
				return 9
			case 1:
				return -9
			}
		case 0x4:
			bits1 = ms.Peekbits(1)
			if bits1&0x01 == 1 {
				_ = ms.Getbits(1)
				return ternary((bits1&0x01) == 1, -10, 10)
			}
			// guaranteed to have 3 bits now
			bits3 = ms.Getbits(3)
			switch bits3 {
			case 2:
			case 3:
				return ternary((bits3&0x01) == 1, -11, 11)
			case 0:
			case 1:
				return ternary((bits3&0x01) == 1, -12, 12)
			default:
				panic("read_mb_mvm: bad 3 bit code (1)");
			}
				panic("read_mb_mvm: bad 3 bit code (2)");
		case 0x3:
			// guaranteed to have 3 bits now
			bits3 = ms.Getbits(3)
			switch bits3 {
			case 6:
			case 7:			
				return ternary((bits3&0x01) == 1, -13, 13)
			case 4:
			case 5:			
				return ternary((bits3&0x01) == 1, -14, 14)
			case 2:
			case 3:			
				return ternary((bits3&0x01) == 1, -15, 15)
			case 0:
			case 1:			
				return ternary((bits3&0x01) == 1, -16, 16)
			default:
				panic("read_mb_mvm: bad 3 bit code (2)");
			}
		default:
			panic("read_mb_mvm: bad 4 bit code (2)");
		}
	default:
		fmt.Printf("bits4a=0x%x\n", bits4a)
		panic("read_mb_mvm: bad 4 bit code (1)")
	}
	panic("read_mb_mvm")
}


func (ms *MpegState) SetYCbCr(mbt *MacroBlockHeader, y0, y1, y2, y3, cb, cr uint32) {
var ternary = func(c bool, a, b bool) bool {
	if (c) {
		return a
	} else {
		return b
	}
}

	mbt.mbt_blockv[0] = ternary((y0 == 1), true, false)
	mbt.mbt_blockv[1] = ternary((y1 == 1), true, false)
	mbt.mbt_blockv[2] = ternary((y2 == 1), true, false)
	mbt.mbt_blockv[3] = ternary((y3 == 1), true, false)
	mbt.mbt_blockv[4] = ternary((cb == 1), true, false)
	mbt.mbt_blockv[5] = ternary((cr == 1), true, false)
	fmt.Printf("iso.SetYCbCr: y0=%v, y1=%v, y2=%v, y3=%v, cb=%v, cr=%v\n",
		mbt.mbt_blockv[0], mbt.mbt_blockv[1], mbt.mbt_blockv[2], mbt.mbt_blockv[3], mbt.mbt_blockv[4], mbt.mbt_blockv[5])
}

func (ms *MpegState) ReadYCbCr() (y0 uint32, y1 uint32, y2 uint32, y3 uint32, cb uint32, cr uint32) {
var bits4a	uint32
var bits4b	uint32
var bits3	uint32
var bits2	uint32
var bits1	bool

	// there is only one 3 bit code, all others are 4 or more
	bits3 = ms.Peekbits(3)
	if bits3 == 7 {
		_ = ms.Getbits(3) 
		y0, y1, y2, y3, cb, cr = 1, 1, 1, 1, 0, 0
	}
	// guaranteed to have 4 bits now
	bits4a = ms.Getbits(4)
	switch bits4a {
	// pure 4 bit codes
	case 0xD:
		y0, y1, y2, y3, cb, cr = 0, 0, 0, 1, 0, 0
	case 0xC:
		y0, y1, y2, y3, cb, cr = 0, 0, 1, 0, 0, 0
	case 0xB:
		y0, y1, y2, y3, cb, cr = 0, 1, 0, 0, 0, 0
	case 0xA:
		y0, y1, y2, y3, cb, cr = 1, 0, 0, 0, 0, 0
	// pure 5 bit codes
	case 0x5:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 0, 0, 1
		} else {
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 1, 0, 1
		}
	case 0x4:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 0, 1, 0
		} else {
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 1, 1, 0
		}
	case 0x9:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 0, 0
		} else {
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 0, 0, 0
		}
	case 0x7:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 1, 0, 0
		} else {
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 1, 0, 0
		}
	case 0x8:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 1, 0, 0
		} else {
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 0, 0, 0
		}
	case 0x6:
		bits1 = ms.Rub()
		if bits1 {
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 1, 0, 0
		} else {
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 0, 0, 0
		}
	// 6 bit codes
	case 0x3:
		// guaranteed 2 bits
		bits2 = ms.Getbits(2)
		switch bits2 {
		case 1:
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 0, 1, 1
		case 3:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 0, 0
		case 2:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 1, 0, 0
		case 0:
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 1, 1, 1
		default:
			panic("read_mb_coded_block_pattern: bad 2 bit code after 0b0011");
		}
	// 7 bit codes
	case 0x2:
		// guaranteed 3 bits
		bits3 = ms.Getbits(3)
		switch bits3 {
		case 0x7:
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 1, 0, 1
		case 0x3:
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 1, 1, 0
		case 0x6:
			y0, y1, y2, y3, cb, cr = 0, 0, 1, 0, 0, 1
		case 0x2:
			y0, y1, y2, y3, cb, cr = 0, 0, 1, 0, 1, 0
		case 0x5:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 0, 0, 1
		case 0x1:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 0, 1, 0
		case 0x4:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 0, 0, 1
		case 0x0:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 0, 1, 0
		default:
			panic("read_mb_coded_block_pattern: bad 3 bit code after 0b0010");
		}
	// 8 bit codes
	case 0x1:
		// guaranteed 4 bits
		bits4b = ms.Getbits(4)
		switch bits4b {
		case 0xF:
			y0, y1, y2, y3, cb, cr = 0, 0, 0, 1, 1, 1
		case 0xE:
			y0, y1, y2, y3, cb, cr = 0, 0, 1, 0, 1, 1
		case 0xB:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 0, 1
		case 0x7:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 1, 0
		case 0x3:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 1, 1
		case 0xD:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 0, 1, 1
		case 0x9:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 1, 0, 1
		case 0x5:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 1, 1, 0
		case 0x1:
			y0, y1, y2, y3, cb, cr = 0, 1, 0, 1, 1, 1
		case 0xC:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 0, 0, 1
		case 0x8:
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 0, 0, 1
		case 0x4:
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 0, 1, 0
		case 0x0:
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 0, 1, 1
		case 0xA:
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 0, 0, 1
		case 0x6:
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 0, 1, 0
		case 0x2:
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 0, 1, 1
		default:
			panic("read_mb_coded_block_pattern: bad 4 bit code after 0b0001");
		}
	// 8-9 bits codes
	case 0x0:
		// at least 4 bits, sometimes 5 next
		bits4b = ms.Getbits(4)
		switch (bits4b) {
		case 0xF:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 0, 1
		case 0xD:
			y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 1, 0
		case 0x1:
			bits1 = ms.Rub()
			if bits1 {
				y0, y1, y2, y3, cb, cr = 0, 1, 1, 0, 1, 1
			} else {
				y0, y1, y2, y3, cb, cr = 1, 0, 0, 1, 1, 1
			}
		case 0xE:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 1, 0, 1
		case 0xC:
			y0, y1, y2, y3, cb, cr = 1, 0, 0, 1, 1, 0
		case 0xA:
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 1, 0, 1
		case 0x6:
			y0, y1, y2, y3, cb, cr = 1, 0, 1, 1, 1, 0
		case 0x3:
			bits1 = ms.Rub()
			if bits1 {
				y0, y1, y2, y3, cb, cr = 0, 1, 1, 1, 1, 1
			} else {
				y0, y1, y2, y3, cb, cr = 1, 0, 1, 1, 1, 1
			}
		case 0x9:
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 1, 0, 1
		case 0x5:
			y0, y1, y2, y3, cb, cr = 1, 1, 0, 1, 1, 0
		case 0x2:
			bits1 = ms.Rub()
			if bits1 {
				y0, y1, y2, y3, cb, cr = 1, 1, 0, 1, 1, 1
			} else {
				y0, y1, y2, y3, cb, cr = 1, 1, 1, 0, 1, 1
			}
		case 0x8:
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 0, 0, 1
		case 0x4:
			y0, y1, y2, y3, cb, cr = 1, 1, 1, 0, 1, 0
		default:
			panic("read_mb_coded_block_pattern: bad 4 bit code after 0b0000")
		}
	default:
		panic("read_mb_coded_block_pattern: bad 4 bit code")
	}
	return y0, y1, y2, y3, cb, cr
}

// read DC difference coding
func (ms *MpegState) ReadDCDC(size int32) int32 {
var bits	uint32
var bits1	uint32
var value	int32 = 1

	if size == 0 {
		return 0
	}
	if (size > 1) {
		bits = ms.Getbits(uint(size-1))
	}
	bits1 = ms.Getbits(1)
	if bits1 == 1 {
		value = int32(^value) * -1
	} else {
		value = int32(bits)
	}
	return value
}

/*
// Each byte consists of a value|length pair
private static final short[] dct_dc_size_luminance = {
  0x12, 0x12, 0x12, 0x12, 0x22, 0x22, 0x22, 0x22,
  0x03, 0x03, 0x33, 0x33, 0x43, 0x43, 0x54, 0x00
};
 
private static final short[] dct_dc_size_luminance1 = {
  0x65, 0x65, 0x65, 0x65, 0x76, 0x76, 0x87, 0x00
}; 
 
public int decodeDCTDCSizeLuminance(InputBitStream input) throws IOException {
  int index = input.nextBits(7);
  int value = dct_dc_size_luminance[index >> 3];
 
  if (value == 0)
    value = dct_dc_size_luminance1[index & 0x07];
 
  input.Getbits(value & 0xf);
 
  return value >> 4
*/

// Read Macro Block DCT DC code Y (luminance)
// first we read the size and then we get the read difference coded value
func (ms *MpegState) ReadMBDCTDCY() (ret int32) {
var		bits1 uint32
var		bits2 uint32
var		cnt1s int32 = 2;
var		size int32

	bits2 = ms.Getbits(2)
	switch bits2 {
	case 0:
		size = 1
	case 1:
		size = 2
	case 2:
		bits1 = ms.Getbits(1)
		if (bits1 == 0) {
			size = 0
		} else {
			size = 3
		}
	case 3:
		for {
			bits1 = ms.Peekbits(1)
			if cnt1s >= 6 || bits1 == 0 { // 7 or 8?
				break
			}
			_ = ms.Getbits(1)
			cnt1s++
		}
		size = cnt1s + 2
	}
	if (size > 0) {
		ret = ms.ReadDCDC(size)
	}
	fmt.Printf("iso.ReadMBDCTDCY: %d\n", ret)
	return
}


// Read Macro Block DCT DC code Cr or Cb (chroma)
// first we read the size and then we read the DC difference coded value
func (ms *MpegState) ReadMBDCTDCC() (ret int32) {
var		bits1 uint32
var		bits2 uint32
var		cnt1s int32 = 2;
var		size int32

	bits2 = ms.Getbits(2)
	switch bits2 {
	case 0:
		size = 0
	case 1:
		size = 1
	case 2:
		size = 2
	case 3:
		for {
			bits1 = ms.Peekbits(1)
			if cnt1s >= 7 || bits1 == 0 { // 7 or 8?
				break
			}
			_ = ms.Getbits(1)
			cnt1s++
		}
		size = cnt1s + 1
	}
	if size > 0 {
		ret = ms.ReadDCDC(size)
	}
	fmt.Printf("iso.ReadMBDCTDCC: %d\n", ret)
	return
}

/*
func (bs *bitstream) ReadRLvlc() {
var		bits2 uint32

	bits2 = bs.Getbits(2)
	switch bits2 {
	case 0:
	case 1:
	case 2:
	case 3:
	}
}
*/

// not needed in Go ?
func mbt_init(mbt *MacroBlockHeader, pt PictureType) {

	mbt.mbt_pt = pt
	mbt.mbt_ai = 0
	mbt.mbt_qs = 0
	
	mbt.mbt_in = false
	mbt.mbt_pa = false
	mbt.mbt_mb = false
	mbt.mbt_mf = false
	mbt.mbt_qf = false
	
	mbt.mbt_mfhp = 0
	mbt.mbt_mfhr = 0
	mbt.mbt_mfvp = 0
	mbt.mbt_mfvr = 0
	mbt.mbt_mbhp = 0
	mbt.mbt_mbhr = 0
	mbt.mbt_mbvp = 0
	mbt.mbt_mbvr = 0
	for i := 0; i < 6; i++ {
		mbt.mbt_blockv[i] = false
	}
}

/*
            mVlc.decodeDCTCoeff(mInput, true, runLevel);

		    run = runLevel.run;
	    	mDctZigzag[run] = runLevel.level;
        }

        if (mPictureCodingType != Picture.D_TYPE) {
            while (mInput.nextBits(2) != 0x2) {
                // dctCoeffNext
            	mVlc.decodeDCTCoeff(mInput, false, runLevel);

                run += runLevel.run + 1;
                mDctZigzag[run] = runLevel.level;
            }
*/

func (ms *MpegState) ReadMacroBlock(mbh *MacroBlockHeader, i int) {
	fmt.Printf("iso.ReadMacroBlock: i=%d\n", i)
	if mbh.mbt_in {
		switch i {
		case 0, 1, 2, 3:
			ms.ReadMBDCTDCY()
		case 4, 5:
			ms.ReadMBDCTDCC()
		}
		if ms.Peekbits(2) == EOB {
			_ = ms.Getbits(2)
			return
		} else {
			fmt.Printf("iso.ReadMacroBlock no EOB, bits2=0x%x\n", ms.Peekbits(2))
		}
		fmt.Printf("iso.ReadMacroBlock: getting coef\n")
		run, level := ms.DecodeDCTCoeff(true)
		fmt.Printf("iso.ReadMacroBlock: first run=%d, level=%d\n", run, level)


		for ms.Peekbits(2) != EOB {
			run, level := ms.DecodeDCTCoeff(false)
			fmt.Printf("iso.ReadMacroBlock: run=%d, level=%d\n", run, level)
		}
		fmt.Printf("iso.ReadMacroBlock: EOB\n")
		ms.Getbits(2)
	}
} 



func (ms *MpegState) ReadMacroBlocks(sh *SequenceHeader, gh *GroupHeader, ph *PictureHeader, sl *SliceHeader) {
var bits11	uint32
var mbh		MacroBlockHeader
var mfhp	int16

var	mfvp	int16

var	mbhp	int16

var	mbvp	int16

var	mfhr	int16
var	mfvr	int16
var	mbhr	int16
var	mbvr	int16


	fmt.Printf("bitstream.ReadMacroBlocks\n")
	mbt_init(&mbh, ph.ph_picture_type)

	for {
		bits11 = ms.Peekbits(11) // check for macro block stuffing
		if bits11 == 0xF {
			_ = ms.Getbits(11)
		} else {
			break
		}
	}
	
	for {
		bits11 = ms.Peekbits(11) // check for macro block escape
		if bits11 == 0x8 {
			_ = ms.Getbits(11)
			mbh.mbt_ai += 33
		} else {
			break
		}
	}
	mbh.mbt_ai = ms.ReadMBAI()	// get address increment
	fmt.Printf("MBAI=%d\n", mbh.mbt_ai)
	in, pa, mb, mf, qf := ms.ReadMBType(ph.ph_picture_type)
	ms.SetMBT(&mbh, in, pa, mb, mf, qf)

	if mbh.mbt_qf {
		mbh.mbt_qs = ms.Russ(5)
		fmt.Printf("iso.ReadMacroBlocks: q=%d\n", mbh.mbt_qs)
	}
	// read motion vectors if present
	if mbh.mbt_mf {
		mfhp = ms.ReadMBMVM()
		if ph.ph_forw_code != 1 {
			if mfhp != 0 && ph.ph_forw_code > 1 {
				mfhr = ms.Rss(uint(ph.ph_forw_code - 1))
			}
		}
		mfvp = ms.ReadMBMVM()
		if (ph.ph_forw_code != 1) {
			if (mfvp != 0 && ph.ph_forw_code > 1) {
				mfvr = ms.Rss(uint(ph.ph_forw_code - 1))
			}
		}
	}		
	if mbh.mbt_mb {
		mbhp = ms.ReadMBMVM()
		if ph.ph_back_code != 1 {
			if mbhp != 0 && ph.ph_back_code > 1 {
				mbhr = ms.Rss(uint(ph.ph_back_code - 1))
			}
		}
		mbvp = ms.ReadMBMVM()
		if ph.ph_back_code != 1 {
			if mbvp != 0 && ph.ph_back_code > 1 {
				mbvr = ms.Rss(uint(ph.ph_back_code - 1))
			}
		}
	}
	fmt.Printf("iso.ReadMacroBlocks: mfhp=%d, mfhr=%d, mfvp=%d, mfvr=%d, mbhp=%d, mbhr=%d, mbvp=%d, mbvr=%d\n",
		mfhp, mfhr, mfvp, mfvr, mbhp, mbhr, mbvp, mbvr)
	if (mbh.mbt_pa) {
		y0, y1, y2, y3, Cb, Cr := ms.ReadYCbCr()
		ms.SetYCbCr(&mbh, y0, y1, y2, y3, Cb, Cr)
	} else {
		ms.SetYCbCr(&mbh, 1, 1, 1, 1, 1, 1)
	}
	for i, v := range mbh.mbt_blockv {
		if v {
			ms.ReadMacroBlock(&mbh, i)
		}
	}
}


func (ms *MpegState) ReadMPEG1Steam() {
var sh				*SequenceHeader
var gh				*GroupHeader
var ph				*PictureHeader
var slh				*SliceHeader
var start_code		uint32
var zseen			int32
var ul				uint32
var vsc				int32
var uc				byte
var	scf				bool
var	vscf			bool

	defer func() {
		if p := recover(); p != nil {
			if p == "EOF" {
				return
			}
			fmt.Printf("unknown error: %v", p)
			return
		}
	}()

	for {
		ms.GetByteAligned()
		for {
			uc = ms.Ruc()
			fmt.Printf("ruc=0x02%x, zseen=%d, scf=%v\n", uc, zseen, scf)
			if scf {
				start_code = 0x100 | uint32(uc)
				zseen = 0
				scf = false;
				break
			} else {
				if uc == 0 {
					zseen++;
				} else {
					if zseen > 1 && uc == 1 {
						scf = true
					}
					zseen = 0
				}
			}
		}
		fmt.Printf("start code = 0x%X\n", start_code)
		if (start_code == PICTURE_START_CODE || start_code > HIGHEST_SLICE_CODE) && vscf {
			fmt.Printf("%d slices\n", vsc)
			vscf = false
		}
		switch start_code {
		case SEQ_HEADER_CODE:
			fmt.Printf("SEQ_HEADER_CODE\n")
			sh = ms.ReadSeqenceHeader()
			fmt.Printf("    sh_hor_size=%d, sh_ver_size=%d, sh_pel_aspect_ratio=%d, sh_picture_rate=%d, sh_bit_rate=%d\n",
				sh.sh_hor_size, sh.sh_ver_size, sh.sh_pel_aspect_ratio, sh.sh_picture_rate, sh.sh_bit_rate)
		case PICTURE_START_CODE:
			ph = ms.ReadPictureHeader()
			fmt.Printf("PICTURE_START_CODE type=%s\n", pt_str[ph.ph_picture_type])
		case GROUP_START_CODE:
			gh = ms.ReadGroupHeader()
			fmt.Printf("GROUP_START_CODE TC=%02d:%02d:%02d:%02d marker=%v, dff=%v\n",
				gh.gh_tc_hr, gh.gh_tc_min, gh.gh_tc_sec, gh.gh_tc_pic, gh.gh_marker_bit, gh.gh_drop_frame_flag)
		case USER_DATA_START_CODE:
			fmt.Printf("USER_DATA_START_CODE\n")
		case SEQ_ERROR_CODE:
			fmt.Printf("SEQ_ERROR_CODE\n")
		case EXTENSION_START_CODE:
			fmt.Printf("EXTENSION_START_CODE\n")
		case RESERVED_CODE:
			fmt.Printf("RESERVED_CODE\n")
		case SEQ_END_CODE:
			fmt.Printf("SEQ_END_CODE\n")
		case ISO_11172_END_CODE:
			fmt.Printf("ISO_11172_END_CODE\n")
		case PACK_START_CODE:
			fmt.Printf("PACK_START_CODE\n")
		case SYSTEM_HEADER_START_CODE:
			fmt.Printf("SYSTEM_HEADER_START_CODE\n")
		default:
			ul = start_code&uint32(SLICE_MASK)
			if ul == 1 {
				vsc = 1
				vscf = true
			} else
				if ul >= 0x02 && ul <= 0xAF {
					vsc++
				} else {
					fmt.Printf("0x%x, unknown start code\n", start_code)
					panic("main: unkown start code")
				}
			slh = ms.ReadSliceHeader(start_code)
			fmt.Printf("VIDEO SLICE CODE 0x%X row start=%d\n", ul, (ul - 1)*16)

			for { 
				ms.ReadMacroBlocks(sh, gh, ph, slh)
				if ms.Peekbits(23) == 0x0 {
					break
				}
			}
		}
	}
}