package zmqDecoder

import (
	"encoding/binary"
	"io/ioutil"
	"testing"
)

func TestImage(t *testing.T) {
	img := ImageData{Series: 1,
		Frame:     0,
		Hash:      "testhash",
		Shape:     [3]int{1030, 1065, 0},
		Type:      "uint32",
		Encoding:  "lz4<",
		Size:      32593,
		StartTime: 0,
		StopTime:  1,
		RealTime:  1,
	}

	data, err := ioutil.ReadFile("testdata/lz4_32bit_dimage-1.0_00044_00000_ZMQframe00002.raw")
	if err != nil {
		t.Error(err)
	}

	img.DataBlob = data
	_, err = img.deflate()
	if err != nil {
		t.Error(err)
	}
	v0 := binary.LittleEndian.Uint32(img.Data[:4])
	v1 := binary.LittleEndian.Uint32(img.Data[len(img.Data)-4:])

	if v0 != 752 || v1 != 752 {
		t.Error("wrong deflated lz4 value")
	}
	img.WriteTiff("testdata")
}

func TestDecodeImage(t *testing.T) {
	mpm := make([][]byte, 4)
	mpm[0], _ = ioutil.ReadFile("testdata/lz4_32bit_dimage-1.0_00044_00000_ZMQframe00000.raw")
	mpm[1], _ = ioutil.ReadFile("testdata/lz4_32bit_dimage-1.0_00044_00000_ZMQframe00001.raw")
	mpm[2], _ = ioutil.ReadFile("testdata/lz4_32bit_dimage-1.0_00044_00000_ZMQframe00002.raw")
	mpm[3], _ = ioutil.ReadFile("testdata/lz4_32bit_dimage-1.0_00044_00000_ZMQframe00003.raw")
	err := Decode(mpm, "testdata")
	if err != nil {
		t.Error(err)
	}
}
