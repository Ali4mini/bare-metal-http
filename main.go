package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	jobs := make(chan net.Conn, 10)

	port := flag.Int("port", 8080, "port number")
	root := flag.String("root", "./", "root of the server")

	flag.Parse()

	address := fmt.Sprintf("127.0.0.1:%d", *port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("address already taken")
	}
	fmt.Printf("listening on: %v", server.Addr())
	for i := 0; i < 5; i++ {
		go worker(i, jobs, root)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("error in Accept")
		}

		jobs <- conn
	}

}
