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

func ListDatabases() (mcp.Tool, server.ToolHandlerFunc) {

	return listDatabases(), listDatabasesHandler
}

func listDatabases() mcp.Tool {

	return mcp.NewTool("list_databases",

		mcp.WithString("cluster",
			mcp.Required(),
			mcp.Description(CLUSTER_PARAMETER_DESCRIPTION),
		),
		mcp.WithDescription("List all databases in a specific Azure Data Explorer cluster"),
	)
}

func listDatabasesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	clusterName, ok := request.Params.Arguments["cluster"].(string)
	if !ok {
		return nil, errors.New("cluster name missing")
	}

	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Use .show databases command
	dataset, err := client.Mgmt(ctx, "", kql.New(".show databases"))
	if err != nil {
		return nil, err
	}

	databaseNames := []string{}

	// Process the results
	for _, row := range dataset.Tables()[0].Rows() {
		// Access database name by column name
		databaseName, err := row.StringByName("DatabaseName")
		if err != nil {
			return nil, err
		}
		//fmt.Println("Database:", databaseName)
		databaseNames = append(databaseNames, databaseName)
	}

	var result ListDatabasesResponse

	result.Databases = databaseNames

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

type ListDatabasesResponse struct {
	Databases []string `json:"databases"`
}
