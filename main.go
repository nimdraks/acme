package main

import (
	"context"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/config"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/rest"
)

func main() {
	// bind stop channel to context
	ctx := context.Background()

	// start REST server
	config := config.Load("C:\\Users\\가족\\Documents\\vscode\\acme\\acme\\config.json")
	server := rest.New(config)
	server.Listen(ctx.Done())
}
