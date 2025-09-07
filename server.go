package main

import (
	"fmt"
	"log/slog"
	"net"
)

func server(port uint) {
	// Get ip address and port from domain

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		slog.Error("error while dns resolve", slog.String("error", err.Error()))
		return
	}

	tcpLis, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		slog.Error("error while listening")
		return
	}
	defer func() {
		if err := tcpLis.Close(); err != nil {
			slog.Error("error while closing listener", slog.String("error", err.Error()))
		}
	}()

	// Accept & return tcpConn which is a linux file under the hood
	tcpConn, err := tcpLis.AcceptTCP()
	if err != nil {
		slog.Error("error while closing", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := tcpConn.Close(); err != nil {
			slog.Error("closing connection", slog.String("error", err.Error()))
		}
	}()

	conn := newTCPConn(tcpConn)
	defer conn.Close()

	if err := conn.Handle(); err != nil {
		slog.Error("closing connection", slog.String("error", err.Error()))
	}
}
