package common

import (
	"github.com/Azure/azure-kusto-go/azkustodata"
)

func GetClient(endpoint string) (*azkustodata.Client, error) {
	// Create a connection string builder with authentication
	kustoConnectionString := azkustodata.NewConnectionStringBuilder(endpoint).WithDefaultAzureCredential()

	// Initialize the client
	client, err := azkustodata.New(kustoConnectionString)
	if err != nil {
		return nil, err
	}
	return client, nil
}
