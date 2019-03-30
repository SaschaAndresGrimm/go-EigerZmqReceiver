# go-EigerZmqReceiver
golang-based zmq receiver to save EIGER stream data as tiff

## About
go-EigerZmqReceiver demonstrates a simple way to receive DECTRIS EIGER and EIGER2 ZMQ image data and save them locally as .tiff files.
Why go? It's simple to write, platform independent, efficient, and compiled.

## Usage
### I want to compile the source
If you are not familiar with go, have a look at Go's excellent [documentation](https://golang.org/doc/install):
```
go get github.com/SaschaAndresGrimm/go-EigerZmqReceiver
go install
go-EigerZmqReceiver
```
Depends on
- [libzmq](https://github.com/zeromq/libzmq)
- [liblz4](https://lz4.github.io/lz4/)


### I want to use the binary
Check the bin folder if the binary is available for your platform and architecture.
- Yes? Run it.
- No? Compile it, run it, and upload it to the bin collection, thanks!

### Flags
```
Usage of go-EigerZmqReceiver:
  -fpath string
    	File path to store images. If empty no files are stored.
  -ip string
    	ip of EIGER2 DCU (default "10.42.41.10")
  -port int
    	EIGER2 zmq port (default 9999)
```

## Like it?
Great. Feel free to contribute!

## Bug or Feature Requests?
Great. Feel free to share your findings and/or contribute!

## Limitations
- no warranty for performance (which depends on multiple factors)
- <32 bit streaming mode not tested
