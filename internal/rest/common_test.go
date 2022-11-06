package rest

import (
	"context"
	"net"
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

type MockConfig struct {
	mockDSN  string
	mockAddr string
}

func (m *MockConfig) GetDSN() string {
	return m.mockDSN
}

func (m *MockConfig) GetAddress() string {
	return m.mockAddr
}

func startServer(ctx context.Context) (string, *Server, error) {
	// get open port
	address, err := getOpenPort()
	if err != nil {
		return "", nil, err
	}

	// start a server
	server := New(&MockConfig{mockAddr: address, mockDSN: "root:1234@tcp(127.0.0.1:3306)/acme"})
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
