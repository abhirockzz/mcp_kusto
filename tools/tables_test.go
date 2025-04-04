package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/abhirockzz/mcp_kusto/common"
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
			Name: "get_table_schema",
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

	// Unmarshal the content to check for the "tables" key
	var output map[string]any
	if err := json.Unmarshal([]byte(content.Text), &output); err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	// Validate the result contains the "tables" key
	tables, ok := output["tables"].([]any)
	if !ok {
		t.Fatal("Expected 'tables' key in unmarshaled output")
	}

	// Verify the cluster name
	expectedClusterName := clusterName // Use the cluster name from the environment variable
	if output["cluster"] != expectedClusterName {
		t.Fatalf("Expected cluster name '%s', but got '%v'", expectedClusterName, output["cluster"])
	}

	// Verify the database name
	expectedDatabaseName := dbName // Use the database name from the environment variable
	if output["database"] != expectedDatabaseName {
		t.Fatalf("Expected database name '%s', but got '%v'", expectedDatabaseName, output["database"])
	}

	// Verify the exact table name is present
	expectedTableName := tableName
	found := false
	for _, table := range tables {
		if table == expectedTableName {
			found = true
			break
		}
	}
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

	// Create a client and set up the test environment
	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

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
	var schemaResponse map[string]any
	if err := json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &schemaResponse); err != nil {
		t.Fatalf("Failed to unmarshal schema response: %v", err)
	}

	// Check table name
	expectedTableName := tableName
	if schemaResponse["Name"] != expectedTableName {
		t.Fatalf("Expected table name '%s', but got '%v'", expectedTableName, schemaResponse["Name"])
	}

	// Check ordered columns
	orderedColumns, ok := schemaResponse["OrderedColumns"].([]interface{})
	if !ok {
		t.Fatal("Expected 'OrderedColumns' key in schema response")
	}

	// Fetch expected column names from environment variables
	expectedColumnNames := os.Getenv("COLUMN_NAMES")
	if expectedColumnNames == "" {
		t.Fatal("Environment variable COLUMN_NAMES is not set")
	}

	expectedColumns := make(map[string]bool)
	for _, col := range strings.Split(expectedColumnNames, ",") {
		expectedColumns[col] = true
	}

	for _, col := range orderedColumns {
		colMap, ok := col.(map[string]any)
		if !ok {
			t.Fatalf("Column is not a valid object: %v", col)
		}

		colName, ok := colMap["Name"].(string)
		if !ok {
			t.Fatalf("Column does not have a valid 'Name' field: %v", colMap)
		}

		delete(expectedColumns, colName)

	}

	if len(expectedColumns) > 0 {
		t.Fatalf("Expected column names not found: %v", expectedColumns)
	}
}
