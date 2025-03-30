package main

import (
	"io"
	"log"
	"net"
	"os"
)

type Progress struct {
	direction string
	bytes     uint64
}

func handleConn(conn net.Conn, in io.Reader, out io.Writer) {
	defer conn.Close()

	pr := make(chan Progress)

	// this reads from stdin and writes to the connection
	go func() {
		n, err := io.Copy(conn, in)
		if err != nil {
			log.Printf("[%s] ERROR: %v\n", conn.RemoteAddr().String(), err)
		}
		pr <- Progress{"sent to conn", uint64(n)}
	}()

	// this reads from the connection output and writes to stdout
	go func() {
		n, err := io.Copy(out, conn)
		if err != nil {
			log.Printf("[%s] ERROR: %v\n", conn.RemoteAddr().String(), err)
		}
		pr <- Progress{"received from conn", uint64(n)}
	}()

	p := <-pr
	log.Printf("[%s] %s: %d bytes\n", conn.RemoteAddr().String(), p.direction, p.bytes)
	p = <-pr
	log.Printf("[%s] %s: %d bytes\n", conn.RemoteAddr().String(), p.direction, p.bytes)
}

func StartTCPServ(port, protocol string) {
	l, err := net.Listen(protocol, ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	log.Printf("Listening on %s:%s\n", protocol, port)
	conn, err := l.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Established connection from %s\n", conn.RemoteAddr().String())
	handleConn(conn, os.Stdin, os.Stdout)
}

func StartTCPClient(host, port, protocol string) {
	conn, err := net.Dial(protocol, host+port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected to %s\n", host+port)
	handleConn(conn, os.Stdin, os.Stdout)
}
