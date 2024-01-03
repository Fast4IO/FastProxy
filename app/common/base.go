package common

import "net"

func CloseConn(conn *net.Conn) {
	(*conn).Close()

}
