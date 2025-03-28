package main

import (
	"github.com/mark3labs/mcp-go/server"
	"log/slog"
)

type Server struct {
	mcpServer *server.MCPServer
}

func NewServer() *Server {
	s := server.NewMCPServer(
		"mcp-wecombot-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	ms := &Server{
		mcpServer: s,
	}
	return ms
}

func (ms *Server) Serve() {
	if err := server.ServeStdio(ms.mcpServer); err != nil {
		slog.Error("error serving stdio", "error", err)
	}
}
