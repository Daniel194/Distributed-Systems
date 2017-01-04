package main

import (
	"fmt"
	"log"
	"net"
	"encoding/binary"
	"time"
	"math/rand"
)

type slaveNode struct {
	addr      *net.UDPAddr
	timeTicks int64
	delta     int64
}

func main() {

	slaves := []string{"localhost:3333", "localhost:4444", "localhost:5555"}

	go runSlave(slaves[0], "Slave #0")
	go runSlave(slaves[1], "Slave #1")
	go runSlave(slaves[2], "Slave #2")
	go runMaster("localhost:2222", slaves)

	var input string
	fmt.Scanln(&input)
}

func checkError(isFatal bool, err error) {
	if err != nil {
		if isFatal {
			log.Fatal(err)
		}
		fmt.Println(err)
	}
}

func runMaster(address string, slavesList []string) {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(true, err)

	sock, err := net.ListenUDP("udp", serverAddr)
	checkError(true, err)
	defer sock.Close()

	slaveNodes := make([]*slaveNode, len(slavesList))

	for i, slaveAddr := range slavesList {
		tempNodeAddr, err := net.ResolveUDPAddr("udp", slaveAddr)

		if err != nil {
			fmt.Println(err)
			continue
		}

		slaveNodes[i] = &slaveNode{addr: tempNodeAddr}
	}

	// Poll all slave nodes for their time
	masterDelta := time.Now().Unix()
	fmt.Printf("Master before adjustment: %v\n", time.Unix(masterDelta, 0))

	pollSlaves(slaveNodes, sock)

	fmt.Print("\n ================================= \n \n")

	// Compute the algorithm for nodes that have responded
	masterDelta = calculateTime(masterDelta, slaveNodes)
	fmt.Printf("Master after adjustment: %v\n", time.Unix(masterDelta, 0))

	// Send the nodes their new time deltas
	tellTheSlaves(slaveNodes, sock)

}

func pollSlaves(slaves []*slaveNode, masterSock *net.UDPConn) {

	for _, slave := range slaves {
		masterSock.WriteToUDP([]byte("Master requesting slave time"), slave.addr)

		buf := make([]byte, 1024)
		n, _, err := masterSock.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(err)
			return
		}

		ticks, bytes := binary.Varint(buf[:n])

		if bytes <= 0 {
			return
		}

		slave.timeTicks = ticks

	}

}

func calculateTime(now int64, slaves []*slaveNode) int64 {
	totalTime := now
	numResponses := 1
	var masterDelta int64

	for _, slave := range slaves {
		numResponses++
		totalTime += slave.timeTicks
	}

	masterDelta = totalTime / int64(numResponses)

	for _, node := range slaves {
		node.delta = masterDelta - node.timeTicks
	}

	return masterDelta
}

func tellTheSlaves(slaves []*slaveNode, sock *net.UDPConn) {

	for _, node := range slaves {
		writeBuf := make([]byte, 1024)
		bytes := binary.PutVarint(writeBuf, node.delta)

		if bytes <= 0 {
			fmt.Println("Error encoding the node's delta")
			continue
		}

		sock.WriteToUDP(writeBuf, node.addr)
	}
}

func runSlave(address string, name string) {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(true, err)

	sock, err := net.ListenUDP("udp", serverAddr)
	checkError(true, err)
	defer sock.Close()

	var delta int64

	rand.Seed(time.Now().UTC().UnixNano())

	readBuf := make([]byte, 1024)
	n, addr, err := sock.ReadFromUDP(readBuf)

	if err != nil {
		fmt.Println(err)
		return
	}

	randomTime := time.Duration(rand.Intn(600)) * time.Second
	now := time.Now().Add(randomTime).Unix()
	fmt.Printf("%v before adjustment: %v \n", name, time.Unix(now + delta, 0))

	writeBuf := make([]byte, 1024)
	bytes := binary.PutVarint(writeBuf, now)

	if bytes <= 0 {
		fmt.Println("Error encoding time")
		return
	}

	sock.WriteToUDP(writeBuf, addr)

	readBuf = make([]byte, 1024)
	n, addr, err = sock.ReadFromUDP(readBuf)

	if err != nil {
		fmt.Println(err)
		return
	}

	delta, bytes = binary.Varint(readBuf[:n])

	if bytes <= 0 {
		return
	}

	fmt.Printf("%v received adjustment of %v \n", name, delta)
	fmt.Printf("%v after adjustment %v\n", name, time.Unix(now + delta, 0))

}