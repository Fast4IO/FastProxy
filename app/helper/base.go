package helper

import (
	"fastProxy/app/config"
	"fastProxy/app/models"
	"fastProxy/app/services"
	"fmt"
	"io"
	"net"
)

type HttpServer struct {
	Server *net.Listener
}

func Start() {
	s := &HttpServer{
		Server: new(net.Listener),
	}
	var err error
	*s.Server, err = net.Listen("tcp", fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))
	if err != nil {
		panic(err)
	}
	for {
		var conn net.Conn
		conn, err = (*s.Server).Accept()
		if err == nil {
			go callback(conn)

		} else {
			conn.Close()
			break
		}
	}
}

func callback(conn net.Conn) {
	defer conn.Close()
	var err interface{}
	var req *services.HTTPRequest
	req, subClient, err := services.NewHTTPRequest(&conn, 4096, new(services.BaseAuth))

	if err != nil {
		conn.Close()
		return
	}
	if !subClient.Ok {
		conn.Close()
		return
	}
	address := req.Host
	err = OutToTCP(address, &conn, req, subClient)
	return
}

func OutToTCP(address string, inConn *net.Conn, req *services.HTTPRequest, subClient models.User) (err error) {
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
