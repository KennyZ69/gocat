package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func sendFile(conn net.Conn, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	filename := filepath.Base(path)

	n, err := io.Copy(conn, f) // copy the file to the conn
	if err != nil {
		fmt.Printf("Error sending file: %s\n", err)
		return err
	}
	fmt.Printf("Sent %d bytes from file %s\n", n, filename)
	conn.Write([]byte("EOF\n"))
	return nil
}

func receiveFile(r *bufio.Reader, conn net.Conn, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file on remote: %s\n", err)
	}
	defer f.Close()

	var totalBytes int64
	for {
		l, err := r.Peek(10)
		if err != nil {

		}
	}
}
