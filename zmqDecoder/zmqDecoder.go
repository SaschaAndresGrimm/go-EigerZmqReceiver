package zmqDecoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
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

	DataBlob    []byte
	Data        []byte
	ByteSize    int
	ElementSize int
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

//decode zmq header frames which contain the meta data for the current acquisition.
//The header frames are sent after arm.
//If Fpath is given, the files are stored as json format in fpath/goZMQ_05d_header.json.
//Else, output meta data to stdout.
func decodeHeader(multiPartMessage [][]byte) error {
	//todo: write dumpConfig to store flat field, pixel mask, and LUT
	glog.Info("decode header")
	info := message2map(multiPartMessage[0])
	metaData := message2map(multiPartMessage[1])

	//flatfield, pixel mask, and LUT
	if len(multiPartMessage) > 2 {
		glog.Infof("header contains additional data which cannot be dumped")
	}

	//save meta data to json file
	if Fpath != "" {
		fname := fmt.Sprintf("goZMQ_%05.f_header.json", info["series"])
		out := path.Join(Fpath, fname)
		glog.Infof("save header data to %s", out)
		os.MkdirAll(Fpath, os.ModePerm)

		jsonFile, err := os.Create(out)
		if err != nil {
			glog.Errorf("Error creating JSON file: %s", err)
			return err
		}
		defer jsonFile.Close()
		jsonFile.Write(multiPartMessage[1])

	} else {
		//print output to stdout
		for key, value := range metaData {
			fmt.Println("*", key, ":", value)
		}
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

	_, err := imageData.inflate()
	if err != nil {
		return err
	}

	if Fpath != "" {
		err := imageData.WriteTiff(Fpath)
		if err != nil {
			return err
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
func (imageData *ImageData) inflate() (int, error) {

	//calculate the byte size of the deflated data
	elements := imageData.Shape[0] * imageData.Shape[1]
	switch imageData.Type {
	case "uint8":
		imageData.ElementSize = 1
		imageData.ByteSize = 1 * elements
	case "uint16":
		imageData.ElementSize = 2
		imageData.ByteSize = 2 * elements
	case "uint32":
		imageData.ElementSize = 4
		imageData.ByteSize = 4 * elements
	default:
		msg := fmt.Sprintf("image type %s not supported", imageData.Type)
		return 0, errors.New(msg)

	}

	imageData.Data = make([]byte, imageData.ByteSize)

	//apply decompression algorithm, accordingly
	switch {
	case strings.HasPrefix(imageData.Encoding, "lz4"):
		return imageData.readLZ4()
	case strings.HasPrefix(imageData.Encoding, "bs"):
		return imageData.readBSLZ4()
	default:
		msg := fmt.Sprintf("encoding %s not supported", imageData.Encoding)
		return 0, errors.New(msg)
	}

}

//readLZ4 deflates a lz4 compressed byte array to an uncompressed byte array.
func (imageData *ImageData) readLZ4() (int, error) {
	return lz4.UncompressBlock(imageData.DataBlob, imageData.Data)
}
