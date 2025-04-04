package main

import (
	"fmt"

	"github.com/abhirockzz/mcp_kusto/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {

	s := server.NewMCPServer(
		"Kusto MCP server",
		"0.0.5",
		server.WithLogging(),
	)

	s.AddTool(tools.ListDatabases())
	s.AddTool(tools.ListTables())
	s.AddTool(tools.GetTableSchema())
	s.AddTool(tools.ExecuteQuery())

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
