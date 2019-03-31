package bitshuffle

// #include "bitshuffle.h"
import "C"

import (
	"unsafe"
)

//BshufDecompressLz4 inflates a bslz4 compressed byte array to an uncompressed byte array.
// Returns the number of consumed bytes in input buffer or negative int if error happened.
func BshufDecompressLz4(in, out unsafe.Pointer, size, elementSize, blockSize int, ) int {
	ret := C.bshuf_decompress_lz4(in,
		out,
		C.ulonglong(size),
		C.ulonglong(elementSize),
		C.ulonglong(blockSize))

	return int(ret)

}
