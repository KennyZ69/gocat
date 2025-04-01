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

func handleConn(conn net.Conn, in io.Reader, out io.Writer, bidir bool) {
	defer conn.Close()

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
			output := executeCmd(command)
			fmt.Printf("Command output: \n%s\n", output)
			conn.Write([]byte(output + "EOF\n"))
			continue
		}

		if bidir {
			fmt.Println("BIDIRECTIONAL COMMUNICATION ENABLED")
			fmt.Printf("%s>", conn.RemoteAddr().String())
			s := bufio.NewScanner(in)
			if s.Scan() {
				conn.Write([]byte(s.Text() + "EOF\n"))
			}
		} else {
			fmt.Printf("Received: %s", msg)
			conn.Write([]byte("\n"))
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

	s := bufio.NewScanner(os.Stdin)
	r := bufio.NewReader(conn)
	for {
		fmt.Printf("%s> ", conn.RemoteAddr().String())
		if !s.Scan() {
			fmt.Println("Error reading from stdin")
			break
		}

		in := s.Text()
		conn.Write([]byte(in + "\n"))
		if in == "exit" {
			fmt.Println("Exiting gracefully ...")
			break
		}

		var resp string
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				// fmt.Println("Server closed connection")
				break
			}
			if line == "EOF\n" {
				break
			}
			resp += line
		}

		fmt.Print(resp)
	}
}
