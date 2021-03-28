package zmqDecoder

//mostly taken from https://github.com/golang/image/blob/master/tiff/

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/golang/glog"
)

var enc = binary.LittleEndian

type ifdEntry struct {
	tag      int
	datatype int
	data     []uint32
}

func (e ifdEntry) putData(p []byte) {
	for _, d := range e.data {
		switch e.datatype {
		case dtByte, dtASCII:
			p[0] = byte(d)
			p = p[1:]
		case dtShort:
			enc.PutUint16(p, uint16(d))
			p = p[2:]
		case dtLong, dtRational:
			enc.PutUint32(p, uint32(d))
			p = p[4:]
		}
	}
}

type byTag []ifdEntry

func (d byTag) Len() int           { return len(d) }
func (d byTag) Less(i, j int) bool { return d[i].tag < d[j].tag }
func (d byTag) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func writeIFD(w io.Writer, ifdOffset int, d []ifdEntry) error {
	var buf [ifdLen]byte
	// Make space for "pointer area" containing IFD entry data
	// longer than 4 bytes.
	parea := make([]byte, 1024)
	pstart := ifdOffset + ifdLen*len(d) + 6
	var o int // Current offset in parea.

	// The IFD has to be written with the tags in ascending order.
	sort.Sort(byTag(d))

	// Write the number of entries in this IFD.
	if err := binary.Write(w, enc, uint16(len(d))); err != nil {
		return err
	}
	for _, ent := range d {
		enc.PutUint16(buf[0:2], uint16(ent.tag))
		enc.PutUint16(buf[2:4], uint16(ent.datatype))
		count := uint32(len(ent.data))
		if ent.datatype == dtRational {
			count /= 2
		}
		enc.PutUint32(buf[4:8], count)
		datalen := int(count * lengths[ent.datatype])
		if datalen <= 4 {
			ent.putData(buf[8:12])
		} else {
			if (o + datalen) > len(parea) {
				newlen := len(parea) + 1024
				for (o + datalen) > newlen {
					newlen += 1024
				}
				newarea := make([]byte, newlen)
				copy(newarea, parea)
				parea = newarea
			}
			ent.putData(parea[o : o+datalen])
			enc.PutUint32(buf[8:12], uint32(pstart+o))
			o += datalen
		}
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
	}
	// The IFD ends with the offset of the next IFD in the file,
	// or zero if it is the last one (page 14).
	if err := binary.Write(w, enc, uint32(0)); err != nil {
		return err
	}
	_, err := w.Write(parea[:o])
	return err
}

// WriteTiff writes the image imgData data to fname.
func (imgData *ImageData) WriteTiff(fpath string) error {
	fname := fmt.Sprintf("goZMQ_%05d_%05d.tiff", imgData.Series, imgData.Frame)
	out := path.Join(fpath, fname)
	glog.Infof("save image %s", out)
	os.MkdirAll(fpath, os.ModePerm)
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.WriteString(w, leHeader)
	if err != nil {
		return err
	}

	// imageLen is the length of the pixel data in bytes.
	// The offset of the IFD is imageLen + 8 header bytes.
	err = binary.Write(w, enc, uint32(imgData.ByteSize+8))
	if err != nil {
		return err
	}

	_, err = w.Write(imgData.Data)
	if err != nil {
		return err
	}

	photometricInterpretation := uint32(pBlackIsZero)
	samplesPerPixel := uint32(1)
	bitsPerSample := []uint32{}

	switch imgData.Type {
	case "uint32":
		bitsPerSample = []uint32{32}
	case "uint16":
		bitsPerSample = []uint32{16}
	case "uint8":
		bitsPerSample = []uint32{8}
	default:
		msg := fmt.Sprintf("dtype %s not yet implemented", imgData.Type)
		return errors.New(msg)
	}

	ifd := []ifdEntry{
		{tImageWidth, dtShort, []uint32{uint32(imgData.Shape[0])}},
		{tImageLength, dtShort, []uint32{uint32(imgData.Shape[1])}},
		{tCompression, dtShort, []uint32{cNone}},
		{tPhotometricInterpretation, dtShort, []uint32{photometricInterpretation}},
		{tStripOffsets, dtLong, []uint32{8}},
		{tRowsPerStrip, dtShort, []uint32{uint32(imgData.Shape[1])}},
		{tStripByteCounts, dtLong, []uint32{uint32(imgData.ByteSize)}},
		{tXResolution, dtRational, []uint32{1333, 1}},
		{tYResolution, dtRational, []uint32{1333, 1}},
		{tResolutionUnit, dtShort, []uint32{resPerCM}},
		{tBitsPerSample, dtShort, bitsPerSample},
		{tSamplesPerPixel, dtShort, []uint32{samplesPerPixel}},
	}

	return writeIFD(w, imgData.ByteSize+8, ifd)
}
