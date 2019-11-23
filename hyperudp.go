package main

import (
	"fmt"
	arg "github.com/alexflint/go-arg"
	"log"
	"runtime"
	"time"
)

var (
	debug bool
)

//Args stores command line argument values
type Args struct {
	IP         string `arg:"positional,required" help:"Hyperion Server IP"`
	Port       int    `help:"Hyperion FlatBuffer Port"`
	UDPPort    int    `help:"UDP Listen Port"`
	Name       string `help:"Name of this stream un Hyperion"`
	Priority   int32  `help:"Priority of this stream in Hyperion"`
	Duration   int32  `help:"How long a single sent frame is displayed if nothing replaces it (in milliseconds)"`
	UDPTimeout int    `help:"Duration after which, if no packet was received no data will be sent to hyperion (in seconds)"`
	GC         int    `help:"How often the garbage collector should be called (in seconds)"`
	Debug      bool   `help:"Print additional debug output"`
}

//Description provides a short text that tells the user what the program does.
//It is displayed when the -h flag is passed
func (Args) Description() string {
	return "Converts an UDP to a Protobuf stream for use with Hyperion"
}

func main() {
	args := Args{
		Port:       19445,
		UDPPort:    1337,
		Name:       "HyperUDP",
		Priority:   100,
		Duration:   1000,
		UDPTimeout: 2,
		GC:         5,
		Debug:      false,
	}
	arg.MustParse(&args)

	debug = args.Debug
	if debug {
		fmt.Println(args)
	}

	var latestData []byte
	sendchan := make(chan []byte, 1)
	recvchan := make(chan []byte, 1)

	go send(sendchan, args)
	go recv(recvchan, args)

	if debug {
		log.Println("Waiting for first packet")
	}
	latestData = <- recvchan
	if debug {
		log.Println("First packet received - starting loop")
	}

	lastpkg := time.Now()
	lastgc := time.Now()
	for {
		for len(recvchan) > 0 {
			latestData = <-recvchan
			lastpkg = time.Now()
		}
		if len(sendchan) == 0 {
			sendchan <- latestData
		}

		// This fixes that the program stops working after a short time
		if time.Since(lastgc) >= time.Duration(args.GC)*time.Second {
			lastgc = time.Now()
			runtime.GC()
			if debug {
				log.Println("Ran GC")
			}
		}
		if time.Since(lastpkg) >= time.Duration(args.UDPTimeout)*time.Second {
			if debug {
				log.Println("Sleeping until a packet arrives because the timeout was exceeded")
			}
			runtime.GC()
			latestData = <-recvchan // Block goroutine until a new packet is received
			lastpkg = time.Now()
			runtime.GC()
			lastgc = time.Now()
			if debug {
				log.Println("Packet received - resuming")
			}
		}
	}
}
