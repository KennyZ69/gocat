package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const HOST = "127.0.0.1"
const PORT = "4466"
const in1 = "Hello from the first"
const in2 = "Hello from the other side"

func TestTCPServ(t *testing.T) {
	go StartTCPServ(PORT, "tcp", false)
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", PORT))
	assert.Nil(t, err)
	defer conn.Close()

	testMsg := "Testing \n"
	_, err = conn.Write([]byte(testMsg))
	if err != nil {
		t.Fatalf("Failed to write to conn: %v\n", err)
	}
}

func TestTCPHandle(t *testing.T) {
	out := new(bytes.Buffer)
	in := bytes.NewReader([]byte(in1))

	ready := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	l, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)
	defer l.Close()

	addr := l.Addr().String()

	go func() {
		<-ready
		conn, err := net.Dial("tcp", addr)
		assert.Nil(t, err)
		handleConn(conn, in, out, false)
		done <- struct{}{}
	}()

	ready <- struct{}{}
	conn, err := l.Accept()
	assert.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, in1, string(buf[:n]))

	n, err = conn.Write([]byte(in2))
	assert.Nil(t, err)
	err = conn.Close()
	assert.Nil(t, err)
	assert.Equal(t, in2, string(out.Bytes()[:n]))

	done <- struct{}{}
}
