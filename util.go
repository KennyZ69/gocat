package main

import (
	"fmt"
	"os/exec"
)

func usage() {
	fmt.Println("Usage: gocat [options]")
	fmt.Println("Options:")
	fmt.Printf("  -l\tListen for incoming connections\n")
	fmt.Printf("  -u\tUse UDP instead of TCP\n")
	fmt.Printf("  -p\tPort to listen on (default: 5443)\n")
	fmt.Printf("  -h\tHost to listen on\n")
	fmt.Printf("  -bi\tBidirectional transfer of commands\n")
	fmt.Println("Examples:")
	fmt.Println("  gocat -l -p 8080\tListen on port 8080")
	fmt.Println("  gocat -h 127.0.0.1 -p 4444\tConnect to localhost on port 4444")
}

func executeCmd(command string) string {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error executing command: %v\n", err)
	}

	return string(out)
}
