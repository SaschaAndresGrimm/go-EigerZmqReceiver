package bitshuffle

// #include <bitshuffle.h>
import "C"

import (
	"unsafe"
)

//BshufDecompressLz4 deflates a bslz4 compressed byte array to an uncompressed byte array.
func BshufDecompressLz4(in, out unsafe.Pointer, size, elementSize int, blockSize uint32) int {
	ret := C.bshuf_decompress_lz4(in,
		out,
		C.ulong(size),
		C.ulong(elementSize),
		C.ulong(blockSize))

	return int(ret)

}
