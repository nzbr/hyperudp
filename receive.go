package main

import (
	"log"
	"net"
	"strconv"
)

func recv(channel chan []byte, args Args) {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(args.UDPPort))
	if err != nil {
		panic(err)
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Starting UDP listener")

	buff := make([]byte, 1024)
	for {
		n, _, err := sock.ReadFromUDP(buff)
		if err != nil {
			log.Printf("UDP Error: %v\n", err)
		}
		if n%3 == 0 {
			channel <- buff[:n]
			buff = make([]byte, 1024)
		} else {
			log.Printf("Packet length %v is not a multiple of 3", n)
		}
	}
}
