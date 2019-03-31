package bitshuffle

// #include "bitshuffle.h"
import "C"

import (
	"unsafe"
)

//BshufDecompressLz4 inflates a bslz4 compressed byte array to an uncompressed byte array.
// Returns the number of consumed bytes in input buffer or negative int if error happened.
func BshufDecompressLz4(in, out unsafe.Pointer, size, elementSize int, blockSize uint32) int {
	ret := C.bshuf_decompress_lz4(in,
		out,
		C.ulonglong(int(size)),
		C.ulonglong(int(elementSize)),
		C.ulonglong(int(blockSize)))

	return int(ret)

}
