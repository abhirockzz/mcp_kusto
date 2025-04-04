package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-kusto-go/azkustodata/kql"
	"github.com/abhirockzz/mcp_kusto/common"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListTables() (mcp.Tool, server.ToolHandlerFunc) {

	return listTables(), listTablesHandler
}

// listTables returns a tool that lists all tables in a specific Azure Data Explorer database.
func listTables() mcp.Tool {

	return mcp.NewTool("list_tables",

		mcp.WithString("cluster",
			mcp.Required(),
			mcp.Description(CLUSTER_PARAMETER_DESCRIPTION),
		),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database to list tables from."),
		),
		mcp.WithDescription("List all tables in a specific Azure Data Explorer database"),
	)
}

// listTablesHandler handles the request to list all tables in a specific Azure Data Explorer database.
func listTablesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	clusterName, ok := request.Params.Arguments["cluster"].(string)
	if !ok {
		return nil, errors.New("cluster name missing")
	}

	dbName, ok := request.Params.Arguments["database"].(string)
	if !ok {
		return nil, errors.New("database name missing")
	}

	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	dataset, err := client.Mgmt(ctx, dbName, kql.New(".show tables"))
	if err != nil {
		return nil, err
	}

	tableNames := []string{}

	// Process the results
	for _, row := range dataset.Tables()[0].Rows() {
		// Access table name by column name
		tableName, err := row.StringByName("TableName")
		if err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}

	result := map[string]any{
		"cluster":  clusterName,
		"database": dbName,
		"tables":   tableNames,
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// GetTableSchema returns a tool that retrieves the schema of a specific table in an Azure Data Explorer database.
func GetTableSchema() (mcp.Tool, server.ToolHandlerFunc) {

	return getSchema(), getSchemaHandler
}

// getSchema returns a tool that retrieves the schema of a specific table in an Azure Data Explorer database.
func getSchema() mcp.Tool {

	return mcp.NewTool("get_table_schema",

		mcp.WithString("cluster",
			mcp.Required(),
			mcp.Description(CLUSTER_PARAMETER_DESCRIPTION),
		),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database."),
		),

		mcp.WithString("table",
			mcp.Required(),
			mcp.Description("Name of the table to get the schema for."),
		),
		mcp.WithDescription("Get the schema of a specific table in an Azure Data Explorer database"),
	)
}

// getSchemaHandler handles the request to retrieve the schema of a specific table in an Azure Data Explorer database.
func getSchemaHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	clusterName, ok := request.Params.Arguments["cluster"].(string)
	if !ok {
		return nil, errors.New("cluster name missing")
	}

	dbName, ok := request.Params.Arguments["database"].(string)
	if !ok {
		return nil, errors.New("database name missing")
	}

	table, ok := request.Params.Arguments["table"].(string)
	if !ok {
		return nil, errors.New("table name missing")
	}

	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	command := kql.New(".show table ").AddTable(table).AddLiteral(" schema as json")

	//fmt.Println("Command:", command.String())

	dataset, err := client.Mgmt(ctx, dbName, command)
	if err != nil {
		return nil, err
	}

	// Process the schema information
	//fmt.Println("Schema for table", table)
	jsonSchema, err := dataset.Tables()[0].Rows()[0].StringByName("Schema")

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(jsonSchema), nil
}
