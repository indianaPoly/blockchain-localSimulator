package network

import (
	"fmt"
	"net"
)

func StartNode(address string) {
    ln, err := net.Listen("tcp", address)
    if err != nil {
        fmt.Println("Error starting node:", err)
        return
    }
    fmt.Println("Node started on", address)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    buffer := make([]byte, 1024)
    _, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }
    fmt.Println("Received message:", string(buffer))
}
