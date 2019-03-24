package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/SaschaAndresGrimm/go-EigerZmqReceiver/zmqDecoder"

	"github.com/golang/glog"
	zmq "github.com/pebbe/zmq4"
)

var (
	ip    string
	port  int
	fpath string
)

func init() {
	flag.StringVar(&ip, "ip", "10.42.41.10", "ip of EIGER2 DCU")
	flag.IntVar(&port, "port", 9999, "EIGER2 zmq port")
	flag.StringVar(&fpath, "fpath", "", "File path to store images. If empty no files are stored.")
	flag.Set("logtostderr", "true")
	flag.Parse()

}

//Receive opens a zmq pull client to tcp://ip:port
//listen for EIGER zmq frames and process them accordingly
func Receive(ip string, port int) {

	context, _ := zmq.NewContext()
	defer context.Term()

	socket, _ := context.NewSocket(zmq.PULL)
	defer socket.Close()
	host := fmt.Sprintf("tcp://%s:%d", ip, port)
	glog.Info("pull from ", host)
	socket.Connect(host)

	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLIN)
	for {
		polledSockets, _ := poller.Poll(time.Millisecond)
		for _, polled := range polledSockets {
			msg := receiveMultipart(polled.Socket)
			glog.Infof("received %d frames\n", len(msg))
			go zmqDecoder.Decode(msg, fpath)
		}
	}
}

//receiveMultipart zmq frames as concatenated byte array
func receiveMultipart(socket *zmq.Socket) [][]byte {

	multiPartMessage := make([][]byte, 9)
	index := 0
	multiPartMessage[index], _ = socket.RecvBytes(0)
	for more, _ := socket.GetRcvmore(); more == true; {
		index++
		multiPartMessage[index], _ = socket.RecvBytes(0)
		more, _ = socket.GetRcvmore()
	}
	return multiPartMessage[:index+1]
}

//start zmq pull receiver and process puled frames
func main() {
	Receive(ip, port)
}
