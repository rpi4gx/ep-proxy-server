package main

import (
	"fmt"
	"net"
	"strconv"
)

func unidirectionalPipe(a net.Conn, b net.Conn) {
	var buffer = make([]byte, 4096)
	go func() {
		for {
			n, err := a.Read(buffer)
			if err != nil {
				a.Close()
				b.Close()
				return
			}
			_, err = b.Write(buffer[:n])
			if err != nil {
				a.Close()
				b.Close()
				return
			}
		}
	}()
}
func handleConnection(pp *proxyPool, serverConn net.Conn) {
	var host, port, err = pp.getProxyFromPool()
	if err != nil {
		fmt.Println(err)
		serverConn.Close()
		return
	}
	clientConn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("error connecting to proxy %v", err)
		return
	}
	go unidirectionalPipe(serverConn, clientConn)
	go unidirectionalPipe(clientConn, serverConn)
}

func serverStart(port int, pp *proxyPool) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port " + strconv.Itoa(port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		go handleConnection(pp, conn)
	}
}
