package main

import "CloudFileServer/server"

func main() {
	s := server.NewServer(server.DefaultConfig())
	s.Run()
}
