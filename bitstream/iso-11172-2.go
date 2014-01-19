package bitstream	// iso111722 rename on import

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

// import . "leb/mpdm/bitstream"

import "fmt"
//import "flag"
//import "os"
//import "io"

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
const HIGHEST_SLICE_CODE		= 0x1AF
const USER_DATA_START_CODE		= 0x1B2
const SEQ_HEADER_CODE			= 0x1B3
const SEQ_ERROR_CODE			= 0x1B4
const EXTENSION_START_CODE		= 0x1B5
const RESERVED_CODE				= 0x1B6
const SEQ_END_CODE				= 0x1B7
const GROUP_START_CODE			= 0x1B8
const SLICE_MASK				= 0xFF


/* System codes */
const ISO_11172_END_CODE		= 0x1B9		
const PACK_START_CODE			= 0x1BA
const SYSTEM_HEADER_START_CODE	= 0x1BB
const PACKET_START_CODE_MASK	= 0xffffff00
const PACKET_START_CODE_PREFIX	= 0x00000100


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
	gh_ghs				[]*PictureHeader
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

type MPEG1 struct {
	bitstream
	sh			[]*SequenceHeader
}

var next_code	uint32

func (bs *bitstream) ReadSeqenceHeader() *SequenceHeader {
var sh			SequenceHeader

//	b |= rul(sp, &shp->sh_code);
	fmt.Printf("bitstream.ReadSeqenceHeader: start\n")
	sh.sh_code = SEQ_HEADER_CODE
	sh.sh_hor_size = bs.russ(12)
	sh.sh_ver_size = bs.russ(12)
	sh.sh_pel_aspect_ratio = bs.russ(4)
	sh.sh_picture_rate = bs.russ(4)
	sh.sh_bit_rate = bs.ruls(18)
	sh.sh_marker_bit = bs.rub()
	sh.sh_vbv_buffer_size = bs.ruls(10)
	sh.sh_const_params_flag = bs.rub()
	sh.sh_load_intra_quant_matrix = bs.rub()
	if sh.sh_load_intra_quant_matrix {
		fmt.Printf("bitstream.ReadSeqenceHeader: 1skip=%d\n", 8*64)
		bs.skipbits(8*64)
	}
	sh.sh_load_non_intra_quant_matrix = bs.rub()
	if sh.sh_load_non_intra_quant_matrix {
		fmt.Printf("bitstream.ReadSeqenceHeader: 2skip=%d\n", 8*64)
		bs.skipbits(8*64)
	}
	next_code = bs.rul()
	sh.sh_iqmp = new(Block)
	sh.sh_niqmp = new(Block)
	sh.sh_extp = nil
	sh.sh_udp = nil
	fmt.Printf("bitstream.ReadSeqenceHeader: end\n")
	return &sh
}


func (bs *bitstream) ReadGroupHeader() *GroupHeader {
var gh			GroupHeader

	fmt.Printf("bitstream.ReadGroupHeader\n")
//	b |= rul(sp, &shp->sh_code);
	gh.gh_code = GROUP_START_CODE;
	gh.gh_drop_frame_flag = bs.rub()
//	b |= ruls(sp, 25, &ghp->gh_tc);
	gh.gh_tc_hr = bs.rucs(5)
	gh.gh_tc_min = bs.rucs(6)
	gh.gh_marker_bit = bs.rub()
	gh.gh_tc_sec = bs.rucs(6)
	gh.gh_tc_pic = bs.rucs(6)
	gh.gh_closed_gap = bs.rub()
	gh.gh_broken_link = bs.rub()
	gh.gh_extp = nil
	gh.gh_udp = nil
	return &gh
}


func (bs *bitstream) getExt() []byte {
var abit		bool
var	cnt			int
var buf 		[100]byte
var bp			[]byte = buf[:]

	for {
		abit = bs.rub()
		if (abit) {
			bp[cnt] = bs.ruc()
			cnt++
		} else {
			break
		}
		if (cnt >= len(buf)) {
			panic("getExt: out of space")
		}
	}
	return bp
}


func (bs *bitstream) ReadPictureHeader() *PictureHeader {
//var abit		bool
//var	i, cnt		int
var uc			byte
//var buf 		[100]byte
//var bp			[]byte = buf[:]
var ph			PictureHeader


	fmt.Printf("bitstream.ReadPictureHeader\n")
//	b |= rul(sp, &shp->sh_code);
	ph.ph_code = PICTURE_START_CODE
	ph.ph_temporal_ref = bs.russ(10)
	uc = bs.rucs(3)
	ph.ph_picture_type = b3_to_PT[uc]
	ph.ph_vbv_delay = bs.rus()
	if ph.ph_picture_type == pt_ppict || ph.ph_picture_type == pt_bpict {
		ph.ph_full_pell_forw_vec= bs.rub()
		ph.ph_forw_code = bs.rcs(3)
	}
	if ph.ph_picture_type == pt_bpict {
		ph.ph_full_pell_forw_vec = bs.rub()
		ph.ph_forw_code = bs.rcs(3)
	}
	
	ph.ph_eip = bs.getExt()
	ph.ph_eip_count = len(ph.ph_eip)
	ph.ph_ext = nil;
	ph.ph_ud = nil;
	return &ph
}


func (bs *bitstream) ReadSliceHeader(code uint32) *SliceHeader {
//var abit bool
//var i, cnt int
//var buf 		[100]byte
//var bp			[]byte = buf[:]
var sl			SliceHeader
	
	fmt.Printf("bitstream.ReadSliceHeader\n")
	sl.sl_code = code;
	sl.sl_quant_scale = bs.rucs(5)
	sl.sl_eip = bs.getExt()
	sl.sl_eip_count = len(sl.sl_eip)
	return &sl
}


func (bs *bitstream) ReadMBAI() uint32 {
var 			bits4a, bits4b, bits3, bits2, bits1 uint32

	// first try the three most common cases 0b1, 0b011, 0b010
	bits1 = bs.peekbits(1)
	if bits1 == 1 {
		_ = bs.getbits(1)
		return 1
	}
	bits3 = bs.peekbits(3)
	if bits3 == 0x3 {
		_ = bs.getbits(3)
		return 2
	}
	bits3 = bs.peekbits(3)
	if bits3 == 0x2 {
		_ = bs.getbits(3)
		return 3
	}
	// ok, it's not a common case, now we basically peel off 4 bits and then determine
	// how many more bits will need to get the whole code which is 4-11 bits
	bits4a = bs.peekbits(4)
	switch bits4a {
	case 0:
		_ = bs.getbits(4)
		bits4b = bs.peekbits(4)
		switch bits4b {
		case 0:
			panic("read_mb_ai: zero stuffed")
			break;
		case 1:
			panic("read_mb_ai: bad pat")
			break;
		case 3:
			_ = bs.getbits(4)
			bits3 = bs.getbits(3)
			return 26 + (7 - bits3)
		case 4:
			_ = bs.getbits(4)
			bits2 = bs.peekbits(2)
			switch bits2 {
			case 0x2:
			case 0x3:
				_ = bs.getbits(2)
				return 20 + (3 - bits2)
			case 0x00:
			case 0x01:
				_ = bs.getbits(2)
				return 22 + (3 - bits2)
			default:
				panic("read_mb_ai: base 4");
			}
		case 5:
			_ = bs.getbits(4)
			bits2 = bs.getbits(2)
			return 16 + (3 - bits2)
		case 6, 7, 8, 9, 10, 11:
			_ = bs.getbits(4)
			bits4b = bs.getbits(4)
			return 10 + (11 - bits4b)
		default:
			panic("read_mb_ai: bad second 4 bits")
		}
		break
	case 2:
	case 3:
		_ = bs.getbits(4)
		return 4 + (3 - bits4a)
	case 1:
		_ = bs.getbits(4)
		bits1 = bs.getbits(1)
		return 6 + (1 - bits1)
	default:
		panic("read_mb_ai: bad first 4 bits");
	}
	panic("read_mb_ai: bad case");
}


func (mbt *MacroBlockHeader) SetMBT(in, pa, mb, mf, qf uint32) {
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
	mbt.mbt_in = ui32_to_bool(in)
	mbt.mbt_pa = ui32_to_bool(pa)
	mbt.mbt_mb = ui32_to_bool(mb)
	mbt.mbt_mf = ui32_to_bool(mf)
	mbt.mbt_qf = ui32_to_bool(qf)
}


func (bs *bitstream) ReadMBType(mbt *MacroBlockHeader) {
var			bits6 uint32
var			bits5 uint32
var			bits4 uint32
var			bits3 uint32
var			bits2 uint32
var			bits1 uint32
	
	bits1 = bs.peekbits(1)
	if (bits1 == 1) {
		_ = bs.getbits(1)
		switch (mbt.mbt_pt) {
		case pt_ipict:
			mbt.SetMBT(1, 0, 0, 0, 0)
		case pt_ppict:
			mbt.SetMBT(0, 1, 0, 1, 0)
		case pt_dpict:
			mbt.SetMBT(1, 0, 0, 0, 0)
		default:
			fmt.Printf("ReadMBType: mbt.mbt_pt=%d\n", mbt.mbt_pt);
			panic("ReadMBType: bad 1");
		}
		return
	}
	bits2 = bs.peekbits(2)
	switch bits2 {
	case 1:
		_ = bs.getbits(2)
		switch mbt.mbt_pt {
		case pt_ipict:
			mbt.SetMBT(1, 0, 0, 0, 1)
		case pt_ppict:
			mbt.SetMBT(0, 1, 0, 0, 0)
		default:
			panic("read_mb_type: bad 2");
		}
		return
	case 2:
		_ = bs.getbits(2)
		if mbt.mbt_pt == pt_ppict {
			mbt.SetMBT(0, 0, 1, 1, 0)
			return
		} else {
			panic("read_mb_type: bad 3")
		}
	case 3:
		_ = bs.getbits(2)
		if mbt.mbt_pt == pt_ppict {
			mbt.SetMBT(0, 1, 1, 1, 0)
			return
		} else {
			panic("read_mb_type: bad 4")
		}
	}
	bits3 = bs.peekbits(3)
	switch(bits3) {
	case 1:
		_ = bs.getbits(3)
		switch mbt.mbt_pt {
		case pt_ipict:
			mbt.SetMBT(1, 0, 0, 0, 1)
		case pt_ppict:
			mbt.SetMBT(0, 1, 0, 0, 0)
		default:
			panic("read_mb_type: bad 2")
		}
		return
	case 2:
		_ = bs.getbits(3) // ???
		if mbt.mbt_pt == pt_ppict {
			mbt.SetMBT(0, 0, 1, 1, 0)
			return
		} else {
			panic("read_mb_type: bad 3")
		}
	case 3:
		_ = bs.getbits(3) // ???
		if mbt.mbt_pt == pt_ppict {
			mbt.SetMBT(0, 1, 1, 1, 0)
			return
		} else {
			panic("read_mb_type: bad 4")
		}
	}
	bits4 = bs.peekbits(4)
	switch bits4 {
	case 0:
		if mbt.mbt_pt == pt_ppict {	// only 1 5 bit code that starts with 0b0000
			bits5 = bs.peekbits(5)
			if bits5 == 1 {
				bits5 = bs.getbits(5)
				mbt.SetMBT(0, 1, 0, 0, 1)
				return
			}
		}
		bits6 = bs.getbits(6)
		switch(bits6) {
		case 0x01:
			if mbt.mbt_pt == pt_ppict {
				mbt.SetMBT(1, 0, 0, 0, 1)
				return
			} else {
				if mbt.mbt_pt == pt_bpict {
					mbt.SetMBT(1, 0, 0, 0, 1)
					return
				} else {
					panic("read_mb_type: bad 0b000001")
				}
			}
		case 0x02:
			if mbt.mbt_pt != pt_bpict {
				panic("read_mb_type: bad 0b000010")
				mbt.SetMBT(0, 1, 1, 0, 1)
				return
			} else {
				panic("read_mb_type: 0b000010 expected pt_bpict")
			}
		case 0x03:
			if mbt.mbt_pt != pt_bpict {
				panic("read_mb_type: bad 0b000011")
				mbt.SetMBT(0, 1, 0, 1, 1)
				return
			} else {
				panic("read_mb_type: 0b000011 expected pt_bpict")
			}
		default:
			panic("read_mb_type: bad 6 bit code")
		}
	case 0x01:	// all 5 bits codes
		bits5 = bs.getbits(5)
		switch bits5 {
		case 0x2:
			if mbt.mbt_pt == pt_ppict {
				mbt.SetMBT(0, 1, 0, 1, 1)
				return
			} else {
				if mbt.mbt_pt == pt_bpict {
					mbt.SetMBT(0, 1, 1, 1, 1)
					return
				} else {
					panic("read_mb_type: bad 0b00010")
				}
			}
		case 0x03:
			if mbt.mbt_pt == pt_ppict {
				mbt.SetMBT(1, 0, 0, 0, 0)
				return
			} else {
				if mbt.mbt_pt == pt_bpict {
					mbt.SetMBT(1, 0, 0, 0, 0)
					return
				} else {
					panic("read_mb_type: bad 0b00011")
				}
			}
		default:
			panic("read_mb_type: bad 5 bit code")
		}
	case 0x02:
		bits4 = bs.getbits(4)
		if mbt.mbt_pt != pt_ppict {
			panic("read_mb_type: code 0b0010 expected pt_ppict")
		}
		mbt.SetMBT(0, 0, 0, 1, 0)
		return
	case 0x03:
		bits4 = bs.getbits(4)
		if mbt.mbt_pt != pt_ppict {
			panic("read_mb_type: code 0b0010 expected pt_ppict")
		}
		mbt.SetMBT(0, 1, 0, 1, 0)
		return
	default:
		panic("read_mb_type: bad 4 bit code")
	}
	panic("read_mb_type: screw")
	return
}




func (bs *bitstream) ReadMBMVM() int16 {
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

	bits1 = bs.peekbits(1)
	if bits1 == 1 {
		_ = bs.getbits(1)
		return 0
	}
	bits3 = bs.peekbits(3)
	if (bits3&0x6) == 0x010 {
		_ = bs.getbits(3)
		return ternary((bits3&0x1) == 1, -1, 1)
	}
	bits4a = bs.peekbits(4)
	switch(bits4a) {
	case 0x2:
	case 0x3:
		_ = bs.getbits(4)
		return ternary((bits4a&0x01) == 1, -2, 2)
	case 0x1:
		_ = bs.getbits(4)
		bits1 = bs.getbits(1)
		return ternary((bits1&0x01) == 1, -3, 3)
	case 0x0:
		// ???
		_ = bs.getbits(4)
		bits3 = bs.peekbits(3)
		if (bits3&0x6) == 0x6 {
			_ = bs.getbits(3)
			return ternary((bits3&0x01) == 1, -4, 4)
		}
		// guaranteed to have 4 bits now, this is the second set of 4 bits
		bits4b = bs.peekbits(4)
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
			bits2 = bs.getbits(2)
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
			bits1 = bs.peekbits(1)
			if bits1&0x01 == 1 {
				_ = bs.getbits(1)
				return ternary((bits1&0x01) == 1, -10, 10)
			}
			// guaranteed to have 3 bits now
			bits3 = bs.getbits(3)
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
			bits3 = bs.getbits(3)
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


func SetYCbCr(mbt *MacroBlockHeader, y0, y1, y2, y3, cb, cr uint32) {
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
}

func (bs *bitstream) ReadYCbCr (mbt *MacroBlockHeader) {
var bits4a	uint32
var bits4b	uint32
var bits3	uint32
var bits2	uint32
var bits1	bool

	// there is only one 3 bit code, all others are 4 or more
	bits3 = bs.peekbits(3)
	if bits3 == 7 {
		_ = bs.getbits(3)
		SetYCbCr(mbt, 1, 1, 1, 1, 0, 0)
		return
	}
	// guaranteed to have 4 bits now
	bits4a = bs.getbits(4)
	switch bits4a {
	// pure 4 bit codes
	case 0xD:
		SetYCbCr(mbt, 0, 0, 0, 1, 0, 0)
	case 0xC:
		SetYCbCr(mbt, 0, 0, 1, 0, 0, 0)
	case 0xB:
		SetYCbCr(mbt, 0, 1, 0, 0, 0, 0)
	case 0xA:
		SetYCbCr(mbt, 1, 0, 0, 0, 0, 0);
	// pure 5 bit codes
	case 0x5:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 0, 0, 0, 0, 0, 1)
		} else {
			SetYCbCr(mbt, 1, 1, 1, 1, 0, 1)
		}
	case 0x4:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 0, 0, 0, 0, 1, 0)
		} else {
			SetYCbCr(mbt, 1, 1, 1, 1, 1, 0)
		}
	case 0x9:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 0, 1, 1, 0, 0, 0)
		} else {
			SetYCbCr(mbt, 1, 1, 0, 0, 0, 0)
		}
	case 0x7:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 0, 1, 1, 1, 0, 0)
		} else {
			SetYCbCr(mbt, 1, 0, 1, 1, 0, 0)
		}
	case 0x8:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 0, 1, 0, 1, 0, 0)
		} else {
			SetYCbCr(mbt, 1, 0, 1, 0, 0, 0)
		}
	case 0x6:
		bits1 = bs.rub()
		if bits1 {
			SetYCbCr(mbt, 1, 1, 0, 1, 0, 0)
		} else {
			SetYCbCr(mbt, 1, 1, 1, 0, 0, 0)
		}
	// 6 bit codes
	case 0x3:
		// guaranteed 2 bits
		bits2 = bs.getbits(2)
		switch bits2 {
		case 1:
			SetYCbCr(mbt, 0, 0, 0, 0, 1, 1)
		case 3:
			SetYCbCr(mbt, 0, 1, 1, 0, 0, 0)
		case 2:
			SetYCbCr(mbt, 1, 0, 0, 1, 0, 0)
		case 0:
			SetYCbCr(mbt, 1, 1, 1, 1, 1, 1)
		default:
			panic("read_mb_coded_block_pattern: bad 2 bit code after 0b0011");
		}
	// 7 bit codes
	case 0x2:
		// guaranteed 3 bits
		bits3 = bs.getbits(3)
		switch bits3 {
		case 0x7:
			SetYCbCr(mbt, 0, 0, 0, 1, 0, 1)
		case 0x3:
			SetYCbCr(mbt, 0, 0, 0, 1, 1, 0)
		case 0x6:
			SetYCbCr(mbt, 0, 0, 1, 0, 0, 1)
		case 0x2:
			SetYCbCr(mbt, 0, 0, 1, 0, 1, 0)
		case 0x5:
			SetYCbCr(mbt, 0, 1, 0, 0, 0, 1)
		case 0x1:
			SetYCbCr(mbt, 0, 1, 0, 0, 1, 0)
		case 0x4:
			SetYCbCr(mbt, 1, 0, 0, 0, 0, 1)
		case 0x0:
			SetYCbCr(mbt, 1, 0, 0, 0, 1, 0)
		default:
			panic("read_mb_coded_block_pattern: bad 3 bit code after 0b0010");
		}
	// 8 bit codes
	case 0x1:
		// guaranteed 4 bits
		bits4b = bs.getbits(4)
		switch bits4b {
		case 0xF:
			SetYCbCr(mbt, 0, 0, 0, 1, 1, 1)
		case 0xE:
			SetYCbCr(mbt, 0, 0, 1, 0, 1, 1)
		case 0xB:
			SetYCbCr(mbt, 0, 1, 1, 0, 0, 1)
		case 0x7:
			SetYCbCr(mbt, 0, 1, 1, 0, 1, 0)
		case 0x3:
			SetYCbCr(mbt, 0, 1, 1, 0, 1, 1)
		case 0xD:
			SetYCbCr(mbt, 0, 1, 0, 0, 1, 1)
		case 0x9:
			SetYCbCr(mbt, 0, 1, 0, 1, 0, 1)
		case 0x5:
			SetYCbCr(mbt, 0, 1, 0, 1, 1, 0)
		case 0x1:
			SetYCbCr(mbt, 0, 1, 0, 1, 1, 1)
		case 0xC:
			SetYCbCr(mbt, 1, 0, 0, 0, 0, 1)
		case 0x8:
			SetYCbCr(mbt, 1, 0, 1, 0, 0, 1)
		case 0x4:
			SetYCbCr(mbt, 1, 0, 1, 0, 1, 0)
		case 0x0:
			SetYCbCr(mbt, 1, 0, 1, 0, 1, 1)
		case 0xA:
			SetYCbCr(mbt, 1, 1, 0, 0, 0, 1)
		case 0x6:
			SetYCbCr(mbt, 1, 1, 0, 0, 1, 0)
		case 0x2:
			SetYCbCr(mbt, 1, 1, 0, 0, 1, 1)
		default:
			panic("read_mb_coded_block_pattern: bad 4 bit code after 0b0001");
		}
	// 8-9 bits codes
	case 0x0:
		// at least 4 bits, sometimes 5 next
		bits4b = bs.getbits(4)
		switch (bits4b) {
		case 0xF:
			SetYCbCr(mbt, 0, 1, 1, 0, 0, 1)
		case 0xD:
			SetYCbCr(mbt, 0, 1, 1, 0, 1, 0)
		case 0x1:
			bits1 = bs.rub()
			if bits1 {
				SetYCbCr(mbt, 0, 1, 1, 0, 1, 1)
			} else {
				SetYCbCr(mbt, 1, 0, 0, 1, 1, 1)
			}
		case 0xE:
			SetYCbCr(mbt, 1, 0, 0, 1, 0, 1)
		case 0xC:
			SetYCbCr(mbt, 1, 0, 0, 1, 1, 0)
		case 0xA:
			SetYCbCr(mbt, 1, 0, 1, 1, 0, 1)
		case 0x6:
			SetYCbCr(mbt, 1, 0, 1, 1, 1, 0)
		case 0x3:
			bits1 = bs.rub()
			if bits1 {
				SetYCbCr(mbt, 0, 1, 1, 1, 1, 1)
			} else {
				SetYCbCr(mbt, 1, 0, 1, 1, 1, 1)
			}
		case 0x9:
			SetYCbCr(mbt, 1, 1, 0, 1, 0, 1)
		case 0x5:
			SetYCbCr(mbt, 1, 1, 0, 1, 1, 0)
		case 0x2:
			bits1 = bs.rub()
			if bits1 {
				SetYCbCr(mbt, 1, 1, 0, 1, 1, 1)
			} else {
				SetYCbCr(mbt, 1, 1, 1, 0, 1, 1)
			}
		case 0x8:
			SetYCbCr(mbt, 1, 1, 1, 0, 0, 1)
		case 0x4:
			SetYCbCr(mbt, 1, 1, 1, 0, 1, 0)
		default:
			panic("read_mb_coded_block_pattern: bad 4 bit code after 0b0000")
		}
	default:
		panic("read_mb_coded_block_pattern: bad 4 bit code")
	}
	return
}

/*
bool
read_mb_dcY_size(strm_t* sp, ulong* sizp)
{
ulong		bits1 = 0;
ulong		bits2 = 0;
ulong		bits3 = 0;
ulong		cnt1s = 2;
bool		b = false;

	b |= getbits(sp, 2, &bits2);
	switch(bits2) {
	case 0b00: *sizp = 1; return(b);
	case 0b01: *sizp = 2; return(b);
	case 0b10:
		b |= getbits(sp, 1, &bits1);
		if (bits1 == 0) {
			*sizp = 0; return(b);
		} else {
			*sizp = 3; return(b);
		break;
	case 0b11:
	do {
		b |= getbits(sp, &bits1);
		cnt1s++;
	} until (cnt1s >= 7 || bits1 == 0);
	*sizp = cnt1s + 2;
	return(b);
}


bool
read_mb_dcC_size(strm_t* sp, ulong* sizp)
{
ulong		bits2 = 0;
ulong		bits3 = 0;
ulong		cnt1s = 2
bool		b = false;

	b |= getbits(sp, 2, &bits2);
	switch(bits2) {
	case 0b00: *sizp = 0; return(b);
	case 0b01: *sizp = 1; return(b);
	case 0b10: *sizp = 2; return(b);
	case 0b11:
	do {
		b |= getbit(sp, &bits1);
		cnt1s++;
	} until (cnt1s >= 8 || bits1 == 0);
	*sizp = cnt1s;
	return(b);
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

var	mfhr	int16
var	mfvr	int16
var	mbhr	int16
var	mbvr	int16

func (bs *bitstream) ReadMacroBlocks(sh *SequenceHeader, gh *GroupHeader, ph *PictureHeader, sl *SliceHeader) {
var bits11	uint32
var mbt		MacroBlockHeader
var mfhp	int16

var	mfvp	int16

var	mbhp	int16

var	mbvp	int16


	fmt.Printf("bitstream.ReadMacroBlocks\n")
	mbt_init(&mbt, ph.ph_picture_type)

	for {
		bits11 = bs.peekbits(11)
		if bits11 == 7 {
			_ = bs.getbits(11)
		} else {
			break
		}
	}
	
	for {
		bits11 = bs.peekbits(11)
		if bits11 == 8 {
			_ = bs.getbits(11)
			mbt.mbt_ai += 33
		} else {
			break
		}
	}
	mbt.mbt_ai = bs.ReadMBAI()
	fmt.Printf("MBAI=%d\n", mbt.mbt_ai)
	bs.ReadMBType(&mbt)
	if mbt.mbt_qf {
		mbt.mbt_qs = bs.russ(5)
	}
	if mbt.mbt_mf {
		mfhp = bs.ReadMBMVM()
		if ph.ph_forw_code != 1 {
			if mfhp != 0 && ph.ph_forw_code > 1 {
				mfhr = bs.rss(uint(ph.ph_forw_code - 1))
			}
		}
		mfvp = bs.ReadMBMVM()
		if (ph.ph_forw_code != 1) {
			if (mfvp != 0 && ph.ph_forw_code > 1) {
				mfvr = bs.rss(uint(ph.ph_forw_code - 1))
			}
		}
	}		
	if mbt.mbt_mb {
		mbhp = bs.ReadMBMVM()
		if ph.ph_back_code != 1 {
			if mbhp != 0 && ph.ph_back_code > 1 {
				mbhr = bs.rss(uint(ph.ph_back_code - 1))
			}
		}
		mbvp = bs.ReadMBMVM()
		if ph.ph_back_code != 1 {
			if mbvp != 0 && ph.ph_back_code > 1 {
				mbvr = bs.rss(uint(ph.ph_back_code - 1))
			}
		}
	}
	bs.ReadYCbCr(&mbt)
}


func (bs *bitstream) ReadMPEG1Steam() {
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
		bs.get_byte_aligned()
		for {
			uc = bs.ruc()
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
			fmt.Printf("SEQ_HEADER_CODE\n");
			sh = bs.ReadSeqenceHeader()
			fmt.Printf("    sh_hor_size=%d, sh_ver_size=%d, sh_pel_aspect_ratio=%d, sh_picture_rate=%d, sh_bit_rate=%lu\n",
				sh.sh_hor_size, sh.sh_ver_size, sh.sh_pel_aspect_ratio, sh.sh_picture_rate, sh.sh_bit_rate)
		case PICTURE_START_CODE:
			ph = bs.ReadPictureHeader()
			fmt.Printf("PICTURE_START_CODE type=%s\n", pt_str[ph.ph_picture_type])
			break;
		case GROUP_START_CODE:
			gh = bs.ReadGroupHeader()
			fmt.Printf("GROUP_START_CODE TC=%02d:%02d:%02d:%02d marker=%d, dff=%d\n",
				gh.gh_tc_hr, gh.gh_tc_min, gh.gh_tc_sec, gh.gh_tc_pic,
					gh.gh_marker_bit, gh.gh_drop_frame_flag)
			break;
		case USER_DATA_START_CODE:
			fmt.Printf("USER_DATA_START_CODE\n");
			break;
		case SEQ_ERROR_CODE:
			fmt.Printf("SEQ_ERROR_CODE\n");
			break;
		case EXTENSION_START_CODE:
			fmt.Printf("EXTENSION_START_CODE\n");
			break;
		case RESERVED_CODE:
			fmt.Printf("RESERVED_CODE\n");
			break;
		case SEQ_END_CODE:
			fmt.Printf("SEQ_END_CODE\n");
			break;
		case ISO_11172_END_CODE:
			fmt.Printf("ISO_11172_END_CODE\n");
			break;
		case PACK_START_CODE:
			fmt.Printf("PACK_START_CODE\n");
			break;
		case SYSTEM_HEADER_START_CODE:
			fmt.Printf("SYSTEM_HEADER_START_CODE\n");
			break;
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
			slh = bs.ReadSliceHeader(start_code)
			bs.ReadMacroBlocks(sh, gh, ph, slh)
			fmt.Printf("VIDEO SLICE CODE 0x%X row start=%d\n", ul, (ul - 1)*16)
			break
		}
	}
}