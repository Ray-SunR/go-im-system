package main

func main() {
	server := NewServer("localhost", 8888)
	server.Start()
}
