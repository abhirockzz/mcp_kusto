package tools

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"slices"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestListTablesHandler(t *testing.T) {
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

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "list_tables",
			Arguments: map[string]any{
				"cluster":  clusterName,
				"database": dbName,
			},
		},
	}

	result, err := listTablesHandler(ctx, request)
	if err != nil {
		t.Fatalf("listTablesHandler failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	content := result.Content[0].(mcp.TextContent)

	t.Logf("Content: %s", content.Text)

	if content.Text == "" {
		t.Fatal("Expected non-empty content")
	}

	// Unmarshal the content to check for the response struct
	var output ListTablesResponse
	if err := json.Unmarshal([]byte(content.Text), &output); err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	// Validate the result contains the expected cluster name
	expectedClusterName := clusterName
	if output.Cluster != expectedClusterName {
		t.Fatalf("Expected cluster name '%s', but got '%v'", expectedClusterName, output.Cluster)
	}

	// Validate the result contains the expected database name
	expectedDatabaseName := dbName
	if output.Database != expectedDatabaseName {
		t.Fatalf("Expected database name '%s', but got '%v'", expectedDatabaseName, output.Database)
	}

	// Validate the result contains the expected table name
	expectedTableName := tableName
	found := slices.Contains(output.Tables, expectedTableName)
	if !found {
		t.Fatalf("Expected table name '%s' not found in tables", expectedTableName)
	}
}

func TestGetSchemaHandler(t *testing.T) {
	ctx := context.Background()

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

	// columnName := os.Getenv("COLUMN_NAME")
	// if columnName == "" {
	// 	t.Fatal("Environment variable COLUMN_NAME is not set")
	// }

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "get_table_schema",
			Arguments: map[string]any{
				"cluster":  clusterName,
				"database": dbName,
				"table":    tableName,
			},
		},
	}

	// Call the handler
	result, err := getSchemaHandler(ctx, request)
	if err != nil {
		t.Fatalf("Handler returned an error: %v", err)
	}

	// Validate the result
	if result == nil {
		t.Fatal("Result is nil")
	}

	schema := result.Content[0].(mcp.TextContent)
	if schema.Text == "" {
		t.Fatal("Schema text is empty")
	}
	t.Logf("Schema text: %s", schema.Text)

	// Verify the schema response
	var schemaResponse TableSchemaResponse
	if err := json.Unmarshal([]byte(schema.Text), &schemaResponse); err != nil {
		t.Fatalf("Failed to unmarshal schema response: %v", err)
	}

	// Check table name
	if schemaResponse.Name != tableName {
		t.Fatalf("Expected table name '%s', but got '%v'", tableName, schemaResponse.Name)
	}

	// Check ordered columns
	expectedColumnNames := os.Getenv("COLUMN_NAMES")
	if expectedColumnNames == "" {
		t.Fatal("Environment variable COLUMN_NAMES is not set")
	}

	expectedColumns := make(map[string]bool)
	for _, col := range strings.Split(expectedColumnNames, ",") {
		expectedColumns[col] = true
	}

	for _, col := range schemaResponse.OrderedColumns {
		delete(expectedColumns, col.Name)
	}

	if len(expectedColumns) > 0 {
		t.Fatalf("Expected column names not found: %v", expectedColumns)
	}
}
