package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Progress struct {
	direction string
	bytes     uint64
}

func handleConn(conn net.Conn, in io.Reader, out io.Writer, bidir bool) {
	defer conn.Close()

	// pr := make(chan Progress)
	// exit := make(chan struct{})

	// INFO: this would be a great way for bidirectional communication but I guess I'll have to rewrite this
	// this reads from stdin and writes to the connection
	// go func() {
	// 	n, err := io.Copy(conn, in)
	// 	if err != nil {
	// 		log.Printf("[%s] ERROR: %v\n", conn.RemoteAddr().String(), err)
	// 	}
	// 	pr <- Progress{"sent to conn", uint64(n)}
	// 	close(exit)
	// }()
	//
	// // this reads from the connection output and writes to stdout
	// go func() {
	// 	n, err := io.Copy(out, conn)
	// 	if err != nil {
	// 		log.Printf("[%s] ERROR: %v\n", conn.RemoteAddr().String(), err)
	// 	}
	// 	pr <- Progress{"received from conn", uint64(n)}
	// 	close(exit)
	// }()
	//
	// select {
	// case p := <-pr:
	// 	// this is the total of bytes sent and received when connection ends
	// 	log.Printf("[%s] %s: %d bytes\n", conn.RemoteAddr().String(), p.direction, p.bytes)
	// case <-exit:
	// 	log.Printf("[%s] Connection closed\n", conn.RemoteAddr().String())
	// }

	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by peer")
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			return
		}

		if strings.HasPrefix(msg, "cmd: ") {
			command := strings.TrimPrefix(msg, "cmd: ")
			log.Printf("Executing command: %s\n", command)
			out := executeCmd(command)
			conn.Write([]byte(out))
			continue
		}

		if bidir {
			fmt.Println("BIDIRECTIONAL COMMUNICATION ENABLED")
			fmt.Printf("%s>", conn.RemoteAddr().String())
			s := bufio.NewScanner(in)
			if s.Scan() {
				conn.Write([]byte(s.Text() + "\n"))
			}
		} else {
			fmt.Printf("Received: %s", msg)
			continue
		}
	}
}

func StartTCPServ(port, protocol string, bidir bool) {
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
	handleConn(conn, os.Stdin, os.Stdout, bidir)
}

func StartTCPClient(host, port, protocol string, bidir bool) {
	addr := net.JoinHostPort(host, port)
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected to %s\n", addr)

	defer conn.Close()
	// handleConn(conn, os.Stdin, os.Stdout, bidir)

	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", conn.RemoteAddr().String())
		if !s.Scan() {
			break
		}

		in := s.Text()
		conn.Write([]byte(in + "\n"))
		if in == "exit" {
			fmt.Println("Exiting gracefully ...")
			break
		}

		rep, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Server closed connection")
			break
		}

		fmt.Printf("Server output: %s", rep)
	}
}
