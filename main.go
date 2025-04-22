package main

import (
	"flag"
	"fmt"
)

var (
	listen = flag.Bool("l", false, "Listen for incoming connections")
	udp    = flag.Bool("u", false, "Use UDP instead of TCP")
	port   = flag.String("p", "5443", "Port to listen on")
	host   = flag.String("h", "", "Host to listen on")
	bidir  = flag.Bool("bi", false, "Bidirectional transfer of commands")
	help   = flag.Bool("help", false, "Show usage info")
)

func main() {
	fmt.Println("Hello from GOCAT!")
	flag.Parse()

	if *help {
		usage()
		return
	}

	if *udp {
		// do udp here
	} else {
		// tcp here
		if *listen {
			StartTCPServ(*port, "tcp", *bidir)
		} else if *host != "" {
			StartTCPClient(*host, *port, "tcp", *bidir)
		} else {
			usage()
		}
	}
}
