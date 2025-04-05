package tools

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestExecuteQueryHandler(t *testing.T) {
	ctx := context.Background()

	// Fetch values from environment variables
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		t.Fatal("Environment variable CLUSTER_NAME is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		t.Fatal("Environment variable DB_NAME is not set")
	}

	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		t.Fatal("Environment variable TABLE_NAME is not set")
	}

	// query := os.Getenv("QUERY")
	// if query == "" {
	// 	t.Fatal("Environment variable QUERY is not set")
	// }

	query := tableName + " | count"

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "execute_query",
			Arguments: map[string]any{
				"cluster":  clusterName,
				"database": dbName,
				"query":    query,
			},
		},
	}

	result, err := executeQueryHandler(ctx, request)
	if err != nil {
		t.Fatalf("executeQueryHandler failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	content := result.Content[0].(mcp.TextContent)

	// t.Logf("Content: %s", content.Text)

	if content.Text == "" {
		t.Fatal("Expected non-empty content")
	}

	//validate content is not empty
	if content.Text == "" {
		t.Fatal("Expected non-empty content")
	}
}
