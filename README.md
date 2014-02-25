mpeg-decoder
============
*Codebase currently under heavy development*

Simple, expositional mpeg-1 decoder written in pure go; not a production quality decoder.

Soon to be added will be the code to perform the iDCT and pixel reconstruction. A (very) simple video viewer (vogl) based on OpenGL is planned.

Note: the code was transliterated from some (incomplete) C code I wrote a number of years ago. The code itself is not too bad at this point. However, some of the structure names still use old style C naming.

Currently the decoder is able to parse mpeg-1 files (IBP frames) and store all the decoded data structures in memory. Note only the ISO-171172-2 file format is supported. The packetized format is not currently supported. However, a demuxer is planned.

Running the program will parse the entire bitstream. If there are no errors there is no output.

MPEG-1 has very little redundancy. Results with corrupted streams are not defined. *This is not a validating parser*

Having said that, the parser does contain an exception based mechanism that could be used to recover from an error and scan for the next start code such as a video slice start code. While the exception mechanism is in place the restart code is not, nor are all the panics that are needed.

On that subject, checking for errors when making many calls to read a few bits at a time is tedious. The bitstream code doesn't return errors on read although there are internal, non exported functions that do. Instead they throw a panic when they get an error, e.g. EOF after reading the last bit. That "exception" can be caught.


Flags
-----
There are flags to print various data structures along the way.

 -from=0: start at frame #  
  -pbc=false: print block coefficients  
  -phd=false: print headers  
  -pmb=false: print macro blocks  
  -prmb=false: print raw macro blocks  
  -ps=false: print stats  
  -pvs=false: print video slices  
  -rmb=true: read macro blocks  
  -to=9999999: stop at frame #  
  -v=false: verbose; turns on most printing

Packages
--------

**iso11172** - This package is basically a toolkit to parse mpeg streams. The current parser will soon be moved out of this package as it is more of an application.

**bitstream** - This package provides routines to parse a variable length bitstream. There is a routine Getbits() to get 1-32 bits and a routine Peekbits() to peek ahead 1-32 bits. Benchamrked on my 2.5 GHz i7 laptop at about 250 Mbit/sec linear read rate.

Features this package doesn't have
----------------------------------
• MPEG-2: Possible, but much more work as there are so many options and variations. Testing and validation would be a huge issue.

• audio:  Need a demuxer for that as well as the audio support and a way to play it back and sync it to the video. Almost certainly need clock support to do it right. I am kind of a nut for A/V sync so this would be a big deal. Also tough to validate without some quality test files which I don't have, although maybe they can be generated with ffmpeg.

• demuxer: It would be nice to have one as more files would be readable. Should be easy to do with Go Routines. Might attempt this.

• de-interlacing: MPEG-1 is defined to be a progressive only video format, but that's often not the case.

• de-blocking filter  
• sharpening filter  
• post processing  
• gamma support  
• accurate color  
• good performance  
• efficiency  
• support for player controls such fast forward, fast backward, slow motion, …  
• brightness, contrast, and color controls  
• scaling of video to window size  




