package main

import (
	"fmt"
	"io"
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.Channel

		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineUsers[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, fmt.Sprintf("User: %s is online", this.Name))
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineUsers, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, fmt.Sprintf("User: %s is offline", this.Name))
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte((msg)))
}

func (this *User) DoMessage() {
	buff := make([]byte, 50)
	for {
		length, err := this.conn.Read(buff)
		if length == 0 {
			this.Offline()
			break
		} else if err != nil && err != io.EOF {
			fmt.Println("Error: ", err)
			break
		} else {
			msg := string(buff[:length-1])
			if msg == "who" {
				this.server.mapLock.RLock()
				for _, user := range this.server.OnlineUsers {
					onlineMessage := "[" + user.Addr + "]" + user.Name + ": is online \n"
					this.SendMessage(onlineMessage)
				}
				this.server.mapLock.RUnlock()
			} else {
				this.server.BroadCast(this, string(buff[:length-1]))
			}
		}
	}
}
