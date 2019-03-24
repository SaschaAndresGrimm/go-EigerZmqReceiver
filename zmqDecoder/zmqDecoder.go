package zmqDecoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/pierrec/lz4"
)

//Fpath is set in case data should be stored on file system.
//If empty nothing is saved, locally.
var Fpath = ""

//ImageData contains the EIGER zmq image meta data as well as the data blob and
//the uncompressed data
type ImageData struct {
	Series int    `json:"series"`
	Frame  int    `json:"frame"`
	Hash   string `json:"hash"`

	Shape    [3]int `json:"shape"`
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Size     int    `json:"size"`

	StartTime int `json:"start_time"`
	StopTime  int `json:"stop_time"`
	RealTime  int `json:"real_time"`

	DataBlob []byte
	Data     []byte
	ByteSize int
}

//Decode EIGER zmq messages and process according to message type
//Save output to fpath if string is not empty
func Decode(multiPartMessage [][]byte, fpath string) error {
	Fpath = fpath

	headerData := message2map(multiPartMessage[0])
	headerType := headerData["htype"].(string)

	switch {
	case strings.Contains(headerType, "dheader"):
		return decodeHeader(multiPartMessage)
	case strings.Contains(headerType, "dimage"):
		return decodeImage(multiPartMessage)
	case strings.Contains(headerType, "dseries_end"):
		return decodeEndOfSeries(multiPartMessage)
	default:
		msg := fmt.Sprintf("header type %s not recognized", headerType)
		return errors.New(msg)
	}
}

//decode zmq header frames which contian meta data for the current acquisition.
//The ehader frames are sent after arm.
func decodeHeader(multiPartMessage [][]byte) error {
	//todo: write dumpConfig to store meta data, flat field, pixel mask, and LUT
	glog.Info("decode header")
	for key, value := range message2map(multiPartMessage[1]) {
		fmt.Println("-", key, ":", value)
	}
	if len(multiPartMessage) > 2 {
		glog.Infof("header contains additional data which cannot be dumped")
	}
	return nil

}

//decode EIGER zmq image data frames to a ImageData struct
func decodeImage(multiPartMessage [][]byte) error {
	glog.Info("decode image")

	imageData := ImageData{}
	json.Unmarshal(multiPartMessage[0], &imageData)
	json.Unmarshal(multiPartMessage[1], &imageData)
	imageData.DataBlob = multiPartMessage[2]
	json.Unmarshal(multiPartMessage[3], &imageData)

	imageData.deflate()
	if Fpath != "" {
		err := imageData.WriteTiff(Fpath)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

//decodeEndOfSeries message whicha re returned after disarm/series end
func decodeEndOfSeries(multiPartMessage [][]byte) error {
	glog.Infof("decode end of series: %s", multiPartMessage)
	return nil
}

//map EIGER zmq frames to an anonymous interface for quick access to
//the json key value pairs
func message2map(message []byte) map[string]interface{} {
	var decoded interface{}
	json.Unmarshal(message, &decoded)
	return decoded.(map[string]interface{})
}

//deflate applies the correct decompression algorithm to the data blob.
//the deflated data are stored in ImageData.Data.
func (imageData *ImageData) deflate() (int, error) {

	//calculate the byte size of the deflated data
	elements := imageData.Shape[0] * imageData.Shape[1]
	switch imageData.Type {
	case "uint16":
		imageData.ByteSize = 2 * elements
	case "uint32":
		imageData.ByteSize = 4 * elements
	default:
		msg := fmt.Sprintf("image type %s not supported", imageData.Type)
		return 0, errors.New(msg)

	}

	imageData.Data = make([]byte, imageData.ByteSize)

	//apply decompression algorithm, accordingly
	switch {
	case strings.HasPrefix(imageData.Encoding, "lz4"):
		return lz4.UncompressBlock(imageData.DataBlob, imageData.Data)
	case strings.HasPrefix(imageData.Encoding, "bs"):
		return readBSLZ4(imageData.DataBlob, imageData.Data)
	default:
		msg := fmt.Sprintf("encoding %s not supported", imageData.Encoding)
		return 0, errors.New(msg)
	}

}

//readLZ4 delates a lz4 compressed byte array to an uncompressed byte array.
func readLZ4(compressed []byte, uncompressed []byte) (int, error) {
	/*
		//eiger data blob needs to have the size prepanded
		prefix := new(bytes.Buffer)
		err := binary.Write(prefix, binary.LittleEndian, int32(len(compressed)))
		if err != nil {
			return 0, err
		}
		fmt.Println(prefix.Bytes())
		compressed = append(prefix.Bytes(), compressed...)
	*/
	return lz4.UncompressBlock(compressed, uncompressed)
}

//readBSLZ4 delates a bslz4 compressed byte array to an uncompressed byte array.
//However, the bitshuffle algorithm is still to be implemented in go. An easy exercise
//for the interested reader.
func readBSLZ4(compressed []byte, uncompressed []byte) (int, error) {
	return 0, errors.New("bs-lz4 not yet implemented")
}
