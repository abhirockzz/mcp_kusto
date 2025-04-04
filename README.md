# MCP server for Azure Data Explorer (Kusto)

This is an implementation of an MCP server for Azure Data Explorer (Kusto) built using its [Go SDK](https://github.com/Azure/azure-kusto-go). It exposes tools for interacting with Azure Data Explorer. You can use this with [VS Code Insiders in Agent mode](https://code.visualstudio.com/blogs/2025/02/24/introducing-copilot-agent-mode) for making data analysis and exploration easier:

1. **list_databases** - Lists all databases in a specific Azure Data Explorer cluster.
   - Parameters:
     - `cluster` (required) - Name of the Azure Data Explorer cluster.

2. **list_tables** - Lists all tables in a specific Azure Data Explorer database.
   - Parameters:
     - `cluster` (required) - Name of the Azure Data Explorer cluster.
     - `database` (required) - Name of the database to list tables from.

3. **get_table_schema** - Gets the schema of a specific table in an Azure Data Explorer database.
   - Parameters:
     - `cluster` (required) - Name of the Azure Data Explorer cluster.
     - `database` (required) - Name of the database.
     - `table` (required) - Name of the table to get the schema for.

4. **execute_query** - Executes a read-only KQL query against a database.
   - Parameters:
     - `cluster` (required) - Name of the Azure Data Explorer cluster.
     - `database` (required) - Name of the database.
     - `query` (required) - The KQL query to execute.

![kusto mcp server in action](mcp_kusto_test.png)

## How to run

```bash
git clone https://github.com/abhirockzz/mcp_kusto
cd mcp_kusto

go build -o mcp_kusto main.go
```

Configure the MCP server. This will differ based on the MCP client/tool you use. For example, with [VS Code](https://code.visualstudio.com/docs/copilot/chat/mcp-servers#:~:text=Choose%20Workspace%20Settings%20to%20create,more%20about%20the%20Configuration%20format.), you can define a `mcp.json` file (inside a `.vscode` folder) as such:

```bash
mkdir -p .vscode

# Define the content for mcp.json
MCP_JSON_CONTENT=$(cat <<EOF
{
  "servers": {
    "Kusto MCP Server": {
      "type": "stdio",
      "command": "$(pwd)/mcp_kusto"
    }
  }
}
EOF
)

# Write the content to mcp.json
echo "$MCP_JSON_CONTENT" > .vscode/mcp.json
```

### Authentication

- The user principal you use should have permissions required for `.show databases`, `.show table`, `.show tables`, and execute queries on the database. Refer to the documentation for [Azure Data Explorer](https://learn.microsoft.com/en-us/kusto/management/security-roles?view=azure-data-explorer) for more details.

- Authentication (Local credentials) - To keep things secure and simple, the MCP server uses [DefaultAzureCredential](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/credential-chains#defaultazurecredential-overview). This approach looks in the environment variables for an application service principal or at locally installed developer tools, such as the Azure CLI, for a set of developer credentials. Either approach can be used to authenticate the MCP server to Azure Data Explorer. For example, just login locally using Azure CLI ([az login](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli)).

You are good to go! Now spin up VS Code Insiders in Agent Mode, or any other MCP tool (like Claude Desktop) and try this out!

## Local dev/testing

Start with [MCP inspector](https://modelcontextprotocol.io/docs/tools/inspector) - `npx @modelcontextprotocol/inspector ./mcp_kusto`
