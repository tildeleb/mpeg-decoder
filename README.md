mpeg-decoder
============

Simple, expositional mpeg-1 decoder written in pure go; not a production quality decoder.

Soon to be added will be the code to perform the iDCT and pixel reconstruction. A (very) simple video viewer based on OpenGL is planned.

Note: the code was hand transliterated from some (incomplete) C code I wrote a number of years ago. The code itself is not too bad at this point. However, some of the structure names are old style C naming.

Currently the decoder is able to parse mpeg-1 files (IBP frames) and store all the decoded data structures in memory. Note only the ISO-171172-2 file format is supported. The packetized format is not currently supported.

Running the program will parse the entire bitstream and if there are no errors there are no results.

MPEG-1 has very little redundancy, the results with corrupted streams are undefined. *This is not a valeting parser*

Having said that, the code does contain an exception based mechanism that could be used to recover from an error and scan for another start code like a video slice start code.

There are flags to print various data structures along the way.

-v prints everything
-ps shows some stats at the end

There are two packages used.

iso11172 - This package is basically a toolkit to parse mpeg streams. The current parser will soon be moved out of this package and it is more of an application.

bitstream - This package provides routines to parse a variable length bitstream. There is a routine Getbits() to get 1-32 bits and a routine Peekbits() to peek ahead 1-32 bits. Benchamrked on my 2.5 GHz i7 laptop at about 250 Mbit/sec linear read rate.

Things you may want that I don't have:
MPEG-2: Possible but much more work as there are lots of options

audio:  Need a demuxed for that as well as the audio support and a way to play and sync it to the video.

de-interlacing: MPEG-1 is suppose to be progressive only but that's often not the case.

deblocking filter
post processing
gamma support
accurate color
good performance
efficiency
support for fast forward, fast backward, slow motion, ...



