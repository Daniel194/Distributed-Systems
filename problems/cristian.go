package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func server() {
	// Listen on a port
	ln, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer ln.Close()

	for {
		// Listen for an incoming connection.
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	// Receive the message
	var msg string

	err := gob.NewDecoder(c).Decode(&msg)

	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	} else {
		fmt.Println("Server Received: ", msg)
	}

	// Send the time
	t := time.Now()
	fmt.Println("Server Sending Time : ", t)

	err = gob.NewEncoder(c).Encode(t)
	if err != nil {
		fmt.Println("Error to send: ", err.Error())
	}

	c.Close()
}

func client() {
	// Connect to the server
	c, err := net.Dial(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)

	if err != nil {
		fmt.Println("Error to connect: ", err.Error())
		os.Exit(1)
	}

	// Send Time
	st := time.Now()

	// Send the message
	msg := "Hello, I want your current time."
	fmt.Println("Client Sending : ", msg)

	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println("Error to send: ", err.Error())
	}

	// Receive the time
	var t time.Time
	err = gob.NewDecoder(c).Decode(&t)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	} else {
		fmt.Println("Client Received Time: ", t)

		// Calculate Round Trip Time
		rtt := time.Since(st)
		fmt.Println("RTT :", rtt)
		rtt /= 2
		fmt.Println("RTT / 2 :", rtt)
		t = t.Add(rtt)

		fmt.Println("Client Received Time + RTT: ", t)
	}

	c.Close()
}

func main() {
	go server()
	go client()

	var input string
	fmt.Scanln(&input)
}