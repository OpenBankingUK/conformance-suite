package main

import (
	"crypto/tls"
	"fmt"
	"log"
)

func main() {
	config := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "golang.cafe:443", &config)
	if err != nil {
		log.Fatalf("unable to establish a conn %v", err)
	}
	defer conn.Close()
	state := conn.ConnectionState()
	fmt.Println(state.Version)
}
