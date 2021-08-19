package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip              string
	Port            int
	OnlineUsers     map[string]*User
	mapLock         sync.RWMutex
	MessageChannnel chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:              ip,
		Port:            port,
		OnlineUsers:     make(map[string]*User),
		MessageChannnel: make(chan string),
	}
	return server
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.MessageChannnel
		this.mapLock.Lock()
		for _, user := range this.OnlineUsers {
			user.Channel <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, message string) {
	message = "[" + user.Addr + "]" + user.Name + ":" + message

	this.MessageChannnel <- message
}

func (this *Server) Handler(conn net.Conn) {
	defer conn.Close()
	user := NewUser(conn, this)
	user.Online()
	user.DoMessage()
}

func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}

	defer listener.Close()

	go this.ListenMessage()
	// accept
	for {
		conn, err := listener.Accept()
		fmt.Println("new connection")
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handlers
		go this.Handler(conn)
	}
	// close listeners
}
