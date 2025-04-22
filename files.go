package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func sendFile(conn net.Conn, path string) error {
	filename := filepath.Base(path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	buf := make([]byte, 4096)

	n, err := f.Read(buf) // read file to its end
	if err != nil {
		f.Close()
		return fmt.Errorf("Error reading file to send: %s\n", err)
	}

	conn.Write(buf[:n])
	fmt.Printf("Sent %d bytes from file %s\n", n, filename)

	// conn.Write([]byte("EOF\n"))
	conn.Write([]byte("<<<EOF>>>\n"))
	f.Close()

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

		content, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			break
			// return fmt.Errorf("Error reading file: %s\n", err)
		}

		if strings.TrimSpace(string(content)) == "<<<EOF>>>\n" {
			fmt.Printf("EOF received\n")
			r.Discard(9)
			break
			// return nil
		}

		n, err := f.Write(content)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			break
			// return fmt.Errorf("Error writing file: %s\n", err)
		}

		totalBytes += int64(n)
	}
	f.Close()
	log.Printf("Received file '%s' [%d bytes]\n", filename, totalBytes)

	return nil
}
