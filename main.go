package main

import "goquickstart/server"

func main() {
	svr := server.NewServer()
	svr.Run()
}
