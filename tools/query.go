package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-kusto-go/azkustodata/kql"
	"github.com/abhirockzz/mcp_kusto/common"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ExecuteQuery() (mcp.Tool, server.ToolHandlerFunc) {

	return executeQuery(), executeQueryHandler
}

func executeQuery() mcp.Tool {

	return mcp.NewTool("execute_query",

		mcp.WithString("cluster",
			mcp.Required(),
			mcp.Description(CLUSTER_PARAMETER_DESCRIPTION),
		),
		mcp.WithString("database",
			mcp.Required(),
			mcp.Description("Name of the database."),
		),

		// mcp.WithString("table",
		// 	mcp.Required(),
		// 	mcp.Description("Name of the table."),
		// ),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The query to execute."),
		),
		mcp.WithDescription("Execute a read-only query. Ask the user for permission before executing the query. It has to be a valid KQL query. Write queries are not allowed. Result truncation is a limit set by default on the result set returned by the query. Kusto limits the number of records returned to the client to 500,000, and the overall data size for those records to 64 MB. When either of these limits is exceeded, the query fails with a partial query failure. Exceeding these limits will generate an exception. Reduce the result set size by modifying the query to only return interesting data. There are several strategies to avoid this. 1/ Use the summarize operator group and aggregate over similar records in the query output. 2/ Potentially sample some columns by using the take_any aggregation function. 3/ Use a take operator to sample the query output. 4/Use the substring function to trim wide free-text columns. 5/ Use the project operator to drop any uninteresting column from the result set."),
	)
}

func executeQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	clusterName, ok := request.Params.Arguments["cluster"].(string)
	if !ok {
		return nil, errors.New("cluster name missing")
	}

	dbName, ok := request.Params.Arguments["database"].(string)
	if !ok {
		return nil, errors.New("database name missing")
	}

	// table, ok := request.Params.Arguments["table"].(string)
	// if !ok {
	// 	return nil, errors.New("table name missing")
	// }

	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("query missing")
	}

	client, err := common.GetClient(fmt.Sprintf(clusterNameFormat, clusterName))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	stmt := kql.New("").AddUnsafe(query)

	queryResponse, err := client.QueryToJson(context.Background(), dbName, stmt)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(queryResponse), nil
}
