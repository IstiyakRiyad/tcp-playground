package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
)

type TCPConn struct {
	tcpConn *net.TCPConn
}

func newTCPConn(conn *net.TCPConn) *TCPConn {
	return &TCPConn{
		tcpConn: conn,
	}
}

func (tc *TCPConn) Close() error {
	if err := tc.tcpConn.Close(); err != nil {
		slog.Error("error while closing tcp conn", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func splitLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0 : i+1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func (tc *TCPConn) Write() error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(splitLine)

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				slog.Error("error while reading line", slog.String("error", err.Error()))
				break
			}
		}
		msg := scanner.Bytes()

		_, err := tc.tcpConn.Write(msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func (tc *TCPConn) Read() error {
	writer := bufio.NewWriter(os.Stdout)
	msg := make([]byte, 1024)

	for {
		_, err := tc.tcpConn.Read(msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if _, err := writer.Write(msg); err != nil {
			return err
		}
		writer.Flush()
		clear(msg)
	}

	return nil
}

func (tc *TCPConn) Handle() error {
	var wg sync.WaitGroup
	var retErr error

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := tc.Read(); err != nil {
			retErr = errors.Join(retErr, err)
		}
		os.Exit(0)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := tc.Write(); err != nil {
			retErr = errors.Join(retErr, err)
		}
		os.Exit(0)
	}()

	wg.Wait()

	return retErr
}
