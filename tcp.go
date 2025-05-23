package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func handleConn(conn net.Conn, in io.Reader, bidir bool) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		command, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by peer")
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			return
		}

		if strings.HasPrefix(command, "msg: ") {
			message := strings.TrimPrefix(command, "msg: ")
			fmt.Printf("Received message: %s\n", message)
			conn.Write([]byte("ACK\nEOF\n"))
			continue
		} else if strings.HasPrefix(command, "file: ") {
			filePath := strings.TrimSpace(strings.TrimPrefix(command, "file: "))
			filename := filepath.Base(filePath)

			// fmt.Printf("Receiving file '%s'. Accept? (y/n)\n", filename)
			// scan := bufio.NewScanner(in)
			// if strings.ToLower(scan.Text()) != "y" {
			// 	conn.Write([]byte("File transfer rejected\nEOF\n"))
			// 	continue
			// }

			conn.Write([]byte("Getting file: " + filename + "\n"))

			if err := receiveFile(r, conn, filename); err != nil {
				log.Printf("Error getting file: %s\n", err)
				conn.Write([]byte("Error receiving file on remote\nEOF\n"))
				continue
			}
			conn.Write([]byte("File received\nEOF\n"))
			continue
		} else {
			log.Printf("Executing command: %s\n", command)
			output := executeCmd(command)
			fmt.Printf("Command output: \n%s\n", output)
			conn.Write([]byte(output + "EOF\n"))
		}

		if bidir {
			fmt.Println("BIDIRECTIONAL COMMUNICATION ENABLED")
			fmt.Printf("%s>", conn.RemoteAddr().String())
			s := bufio.NewScanner(in)
			if s.Scan() {
				conn.Write([]byte(s.Text() + "EOF\n"))
			}
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
	handleConn(conn, os.Stdin, bidir)
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
		} else if strings.HasPrefix(in, "file: ") {
			filePath := strings.TrimSpace(strings.TrimPrefix(in, "file: "))
			if err := sendFile(conn, filePath); err != nil {
				log.Printf("Error sending file: %s\n", err)
			}
			continue
		}

		var resp string
		for {
			line, err := r.ReadString('\n')
			if err != nil {
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
