package net

import (
	"errors"
	"net"
)

func Server(network string, address string) error {
	listen, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		go func() {
			if connErr := handleConn(conn); connErr != nil {
				conn.Close()
			}
		}()
	}

}

func handleConn(conn net.Conn) error {
	for {
		bs := make([]byte, 8)
		n, err := conn.Read(bs)
		if err != nil {
			return err
		}
		if n != len(bs) {
			return errors.New("数据读取不完整")
		}
		res := handleMsg(bs)
		n, err = conn.Write(res)
		if err != nil {
			return err
		}
		if n != len(res) {
			return errors.New("数据写入不完整  ")
		}
	}
}

func handleMsg(req []byte) []byte {
	res := make([]byte, 2*len(req))
	copy(res[:len(req)], req)
	copy(res[len(req):], req)
	return res
}
