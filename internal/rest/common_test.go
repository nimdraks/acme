package rest

import (
	"context"
	"net"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/config"
)

func getOpenPort() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}

	address := listener.Addr().String()
	listener.Close()

	return address, nil
}

func startServer(ctx context.Context) (string, *Server, error) {
	// get open port
	address, err := getOpenPort()
	if err != nil {
		return "", nil, err
	}

	// start a server
	config := config.Config{Address: address, DSN: config.App.DSN}
	server := New(&config)
	go server.Listen(ctx.Done())

	// wait for server to be ready
	dialer := &net.Dialer{}
	for {
		conn, _ := dialer.DialContext(ctx, "tcp", address)
		if conn != nil {
			defer conn.Close()

			return address, server, nil
		}

		select {
		case <-ctx.Done():
			return "", nil, ctx.Err()

		default:
			// try again
		}
	}

	return address, server, nil
}
