package tools

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestListDatabasesHandler(t *testing.T) {
	ctx := context.Background()

	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		t.Fatal("Environment variable CLUSTER_NAME is not set")
	}

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "list_databases",
			Arguments: map[string]interface{}{
				"cluster": clusterName,
			},
		},
	}

	result, err := listDatabasesHandler(ctx, request)
	if err != nil {
		t.Fatalf("listDatabasesHandler failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	content := result.Content[0].(mcp.TextContent)

	t.Logf("Content: %s", content.Text)

	if content.Text == "" {
		t.Fatal("Expected non-empty content")
	}
	// Unmarshal the content to check for the "databases" key
	var output map[string]any
	if err := json.Unmarshal([]byte(content.Text), &output); err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	// Validate the result contains the "databases" key
	if _, ok := output["databases"]; !ok {
		t.Fatal("Expected 'databases' key in unmarshaled output")
	}
}
