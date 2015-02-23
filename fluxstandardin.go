// Copyright (c) 2015 Joshua Scoggins
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgement in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.
//

// reads input from a socket and outputs it to standard out (useful for programs which aren't network aware)
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DrItanium/neuron"
	"log"
	"net"
	"os"
	"syscall"
)

var portNumber = flag.Uint("port", 2000, "the port number to listen on")

func main() {
	terminator := make(chan bool, 1)
	flag.Parse()
	// setup the port
	str := fmt.Sprintf(":%d", *portNumber)
	l, err := net.Listen("tcp", str)
	if err != nil {
		log.Fatal(err)
	}
	neuron.StopRunningOnSignalAndForward(syscall.SIGINT, terminator)
	// this program won't exit conventionally so just terminate if Ctrl-C is pressed
	go func() {
		<-terminator
		l.Close()
		os.Exit(0)
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			a := make([]byte, 80)
			var b bytes.Buffer
			for {
				if !neuron.IsRunning() {
					break
				}
				result, err := c.Read(a)
				if err != nil {
					if err.Error() == "EOF" {
						break
					} else {
						log.Fatal(err)
					}
				}

				if result == 0 {
					break
				} else {
					b.Write(a)
					b.WriteTo(os.Stdout)
				}
			}
			/*
				// Pretty much the callback is something that we shouldn't worry about at this point
					a[0] = 0
					io.Copy(c, bytes.NewReader(a))
			*/
			c.Close()
		}(conn)
	}
}
