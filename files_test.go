package main

import (
	"net"
	"os"
	"testing"
)

func TestFileTransfer(t *testing.T) {
	client, serv := net.Pipe()
	defer client.Close()
	defer serv.Close()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %s\n", err)
	}
	defer os.Remove(tmpFile.Name())

	data := []byte("Testing for file transfer")
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("Failed to write to temp file: %s\n", err)
	}
	tmpFile.Close()

	// receive now
	go func() {
		err := receiveFile(client, serv)
	}()
}
