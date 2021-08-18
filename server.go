package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection established")
	defer func() {
		fmt.Println("Connection cloesd")
		conn.Close()
	}()
	buff := make([]byte, 50)
	for {
		length, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		} else {
			fmt.Println("Read: ", string(buff[:length]))
		}
	}
}

func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	defer listener.Close()
	// accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handlers
		go this.Handler(conn)
	}
	// close listeners
}
