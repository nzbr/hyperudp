package main

import (
	json "encoding/json"
	"fmt"

	arg "github.com/alexflint/go-arg"
)

//Args stores command line argument values
type Args struct {
	IP       string `arg:"positional,required" help:"Hyperion Server IP"`
	Port     int    `help:"Hyperion Protocol Buffer Port"`
	UDPPort  int    `help:"UDP Listen Port"`
	Name     string `help:"Name of this stream in Hyperion"`
	Priority int32  `help:"Priority of this stream in Hyperion"`
	Duration int32  `help:"How long a single sent frame is displayed if nothing replaces it (in milliseconds)"`
}

//Description provides a short text that tells the user what the program does.
//It is displayed when the -h flag is passed
func (Args) Description() string {
	return "Converts an UDP to a Protobuf stream for use with Hyperion"
}

func main() {
	args := Args{
		Port:     19445,
		UDPPort:  1337,
		Name:     "HyperUDP",
		Priority: 100,
		Duration: 1000,
	}
	arg.MustParse(&args)

	data, _ := json.MarshalIndent(args, "", "    ")
	fmt.Println(string(data))

	channel := make(chan []byte, 1)

	go send(channel, args)
	recv(channel, args)
	//recv will never return
}
