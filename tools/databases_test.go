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

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		t.Fatal("Environment variable DB_NAME is not set")
	}

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "list_databases",
			Arguments: map[string]any{
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

	// Unmarshal the content
	var output ListDatabasesResponse
	if err := json.Unmarshal([]byte(content.Text), &output); err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	if len(output.Databases) == 0 {
		t.Fatal("Expected 'databases' key in unmarshaled output with non-empty slice")
	}

	if output.Databases[0] != dbName {
		t.Fatalf("Expected database name %s, got %s", dbName, output.Databases[0])
	}
}

// func createDatabase(clusterName, dbName string) error {

// 	if dbName == "" {
// 		return errors.New("database name cannot be empty")
// 	}

// 	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
// 	if err != nil {
// 		return fmt.Errorf("failed to get client: %v", err)
// 	}

// 	_, err = client.Mgmt(context.Background(), "", kql.New("").AddUnsafe(".create database "+dbName))
// 	if err != nil {
// 		return fmt.Errorf("failed to create database: %v", err)
// 	}
// 	return nil
// }

// func deleteDatabase(clusterName, dbName string) error {

// 	if dbName == "" {
// 		return errors.New("database name cannot be empty")
// 	}

// 	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
// 	if err != nil {
// 		return fmt.Errorf("failed to get client: %v", err)
// 	}

// 	_, err = client.Mgmt(context.Background(), "", kql.New(".drop database ").AddUnsafe(dbName))
// 	if err != nil {
// 		return fmt.Errorf("failed to drop database: %v", err)
// 	}

// 	return nil
// }
