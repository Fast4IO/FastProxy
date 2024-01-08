package services

import (
	"io"
	"net"
)

func OutToTCP(address string, inConn *net.Conn, req *HTTPRequest) (err error) {
	var outConn net.Conn

	outConn, err = net.Dial("tcp", address)
	if err != nil {
		return
	}

	defer outConn.Close()
	if req.Method == "CONNECT" {
		(*inConn).Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		outConn.Write(req.HeadBuf)
	}
	go func() {
		defer outConn.Close()
		defer (*inConn).Close()
		io.Copy(outConn, *inConn)
	}()
	io.Copy(*inConn, outConn)
	return
}
