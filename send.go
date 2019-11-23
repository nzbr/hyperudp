package main

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	hyperion "hyperudp/proto"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func send(channel chan []byte, args Args) {
	var payload []byte

	if debug {
		log.Println("Connecting to Hyperion...")
	}
	sock := connect(args.IP, args.Port)
	if debug {
		log.Println("Connected!")
	}

	for {
		payload = <-channel
		if buff := buildBuffer(payload, args.Priority, args.Duration); buff != nil {
			if err := sendBuffer(sock, buff); err != nil {
				if err == io.EOF {
					log.Println("Disconnected! Trying to reconnect...")
					sock = connect(args.IP, args.Port)
					log.Println("Reconnected successfully!")
				} else {
					panic(err)
				}
			}
		} else {
			log.Println("Could not create buffer from received payload!")
		}
	}
}

func connect(ip string, port int) net.Conn {
	var err error
	var sock net.Conn

	err = fmt.Errorf("Not Initialized")
	start := time.Now()
	next := 10

	for err != nil {
		if debug && time.Since(start).Seconds() >= float64(next) {
			log.Printf("Trying to connect for %v seconds\n", next)
			next += 10
		}
		sock, err = net.Dial("tcp", ip+":"+strconv.Itoa(port))
		if err != nil {
			time.Sleep(time.Millisecond * 100)
		}
	}

	return sock
}

func sendBuffer(socket net.Conn, buffer []byte) error {
	size := len(buffer)

	header := []byte{
		byte(size >> 24 & 0xFF),
		byte(size >> 16 & 0xFF),
		byte(size >> 8 & 0xFF),
		byte(size & 0xFF),
	}

	payload := append(header, buffer...)

	if n, err := socket.Write(payload); n != len(payload) || err != nil {
		if err == nil {
			err = fmt.Errorf("Could not write full payload")
			log.Printf("%v\t%3v/%3v\n", err, n, len(payload))
		}
		return err
	}

	if err := getReply(socket); err != nil {
		return err
	}

	return nil
}

func buildBuffer(udpBytes []byte, priority int32, duration int32) []byte {
	if len(udpBytes)%3 != 0 {
		log.Print("Payload length not divisible by 3!")
		return nil
	}

	width := int32(len(udpBytes) / 3)
	height := int32(1)

	imagerequest := hyperion.ImageRequest{
		Imagedata:   udpBytes,
		Imagewidth:  &width,
		Imageheight: &height,
		Priority:    &priority,
		Duration:    &duration,
	}

	request := hyperion.HyperionRequest{
		Command: hyperion.HyperionRequest_IMAGE.Enum(),
	}
	if err := proto.SetExtension(&request, hyperion.E_ImageRequest_ImageRequest, &imagerequest); err != nil {
		panic(err)
	}

	bytes, err := proto.Marshal(&request)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getReply(sock net.Conn) error {
	header := make([]byte, 4)
	if n, err := sock.Read(header); n != 4 || err != nil {
		if err == nil {
			err = fmt.Errorf("Could not reaed response header")
			log.Printf("%v\t%3v/%3v\n", err, n, 4)
		}
		return err
	}

	size := (int(header[0]) << 24) | (int(header[1]) << 16) | (int(header[2]) << 8) | (int(header[3]))

	data := make([]byte, size)
	if n, err := sock.Read(data); n != size || err != nil {
		if err == nil {
			err = fmt.Errorf("Could not read full response")
			log.Printf("%v\t%3v/%3v\n", err, n, size)
		}
		return err
	}

	reply := &hyperion.HyperionReply{}
	if err := proto.Unmarshal(data, reply); err != nil {
		log.Println(err)
		return err
	}
	if !reply.GetSuccess() {
		fmt.Printf("REMOTE: %v\n", reply.GetError())
	}

	return nil
}
