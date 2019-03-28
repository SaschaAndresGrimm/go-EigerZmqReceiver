package zmqDecoder

// #cgo LDFLAGS: -llz4
// #include "bitshuffle.h"
import "C"
import (
	"encoding/binary"
	"unsafe"
)

func bshufDecompressLZ4(in []byte, out []byte, size uint, elementSize uint, blockSize uint) int {
	ret := C.bshuf_decompress_lz4(unsafe.Pointer(&in[0]), unsafe.Pointer(&out[0]),
		C.ulong(size), C.ulong(elementSize), C.ulong(blockSize))
	return int(ret)
}

//readBSLZ4 deflates a bslz4 compressed byte array to an uncompressed byte array.

func (imgData *ImageData) readBSLZ4() (int, error) {
	//blocksize is big endian uint32 starting at byte 8, divided by element size
	blockSize := binary.BigEndian.Uint32(imgData.DataBlob[8:12]) / uint32(imgData.ElementSize)
	ret := bshufDecompressLZ4(imgData.DataBlob[12:], imgData.Data, uint(imgData.ByteSize),
		uint(imgData.ElementSize), uint(blockSize))
	/*
		    missing something. always getting negative value back which indicates an error.
				    if ret <= 0 {
						msg := fmt.Sprintf("bslz4 decompression error %d", ret)
						fmt.Println(imgData.Data)
						return 0, errors.New(msg)
					}
	*/
	return ret, nil
}
