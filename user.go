package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
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
