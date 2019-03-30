package zmqDecoder

import (
	"io/ioutil"
	"testing"
)

func TestDecodeLZ4_1M(t *testing.T) {
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

func TestDecodeBSLZ4_1M(t *testing.T) {
	mpm := make([][]byte, 4)
	mpm[0], _ = ioutil.ReadFile("testdata/bslz4_32bit_dimage-1.0_00013_00000_ZMQframe00000.raw")
	mpm[1], _ = ioutil.ReadFile("testdata/bslz4_32bit_dimage-1.0_00013_00000_ZMQframe00001.raw")
	mpm[2], _ = ioutil.ReadFile("testdata/bslz4_32bit_dimage-1.0_00013_00000_ZMQframe00002.raw")
	mpm[3], _ = ioutil.ReadFile("testdata/bslz4_32bit_dimage-1.0_00013_00000_ZMQframe00003.raw")
	err := Decode(mpm, "testdata")
	if err != nil {
		t.Error(err)
	}
}
