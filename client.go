package main

import (
	"fmt"
	"log/slog"
	"net"
)

func client(host string, port uint) {
	// laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, 8233))
	// if err != nil {
	// 	slog.Error("error while resolving addr", "error", err.Error())
	// 	return
	// }

	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		slog.Error("error while resolving addr", slog.String("error", err.Error()))
		return
	}

	tcpConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		slog.Error("error while tcp conn", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := tcpConn.Close(); err != nil {
			slog.Error("error while closing conn", slog.String("error", err.Error()))
		}
	}()

	conn := newTCPConn(tcpConn)
	defer conn.Close()

	if err := conn.Handle(); err != nil {
		slog.Error("closing connection", slog.String("error", err.Error()))
	}
}
