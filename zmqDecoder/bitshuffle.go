package zmqDecoder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"github.com/SaschaAndresGrimm/go-EigerZmqReceiver/zmqDecoder/bitshuffle"
)

//readBSLZ4 deflates a bslz4 compressed byte array to an uncompressed byte array.
func (imgData *ImageData) readBSLZ4() (int, error) {
	//blocksize is big endian uint32 starting at byte 8, divided by element size
	//data blob starts at byte 12
	blockSize := int(binary.BigEndian.Uint32(imgData.DataBlob[8:12]) / uint32(imgData.ElementSize))

	ret := bitshuffle.BshufDecompressLz4(unsafe.Pointer(&imgData.DataBlob[12]),
		unsafe.Pointer(&imgData.Data[0]),
		imgData.ByteSize/imgData.ElementSize,
		imgData.ElementSize,
		blockSize)

	//number of bytes consumed in *input* buffer, negative error-code if failed.
	if ret <= 0 {
		msg := fmt.Sprintf("bslz4 decompression error %d", ret)
		fmt.Println(imgData.Data)
		return 0, errors.New(msg)
	}

	return ret, nil

}
