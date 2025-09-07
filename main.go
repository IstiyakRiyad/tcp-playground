package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [host] [port]\n", os.Args[0])
		flag.PrintDefaults()
	}

	var port uint
	flag.UintVar(&port, "l", 0, "(port) only use for server setup")

	flag.Parse()
	args := flag.Args()

	// Client Check
	if len(args) == 2 {
		port, err := strconv.Atoi(args[1])
		if err != nil {
			flag.Usage()
			return
		}

		client(flag.Arg(0), uint(port))
		return
	}

	// Check for server
	if port == 0 {
		flag.Usage()
		return
	}

	server(port)
}
