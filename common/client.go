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

// const clusterNameFormat = "https://%s.kusto.windows.net/"

// func CreateTable(clusterName, dbName, createTableCommand string) error {

// 	endpoint := fmt.Sprintf(clusterNameFormat, clusterName)
// 	fmt.Println("Cluster URL:", endpoint)

// 	// Initialize the client
// 	client, err := GetClient(endpoint)
// 	if err != nil {
// 		return fmt.Errorf("error creating Kusto client: %w", err)
// 	}
// 	defer client.Close()

// 	ctx := context.Background()

// 	_, err = client.Mgmt(ctx, dbName, kql.New("").AddUnsafe(createTableCommand))
// 	if err != nil {
// 		return fmt.Errorf("error executing create table: %w", err)
// 	}

// 	log.Println("table created successfully with command", createTableCommand)
// 	return nil
// }

// func DropTable(clusterName, dbName, tableName string) error {
// 	endpoint := fmt.Sprintf(clusterNameFormat, clusterName)
// 	fmt.Println("Cluster URL:", endpoint)

// 	// Initialize the client
// 	client, err := GetClient(endpoint)
// 	if err != nil {
// 		return fmt.Errorf("error creating Kusto client: %w", err)
// 	}
// 	defer client.Close()

// 	ctx := context.Background()

// 	_, err = client.Mgmt(ctx, dbName, kql.New(". drop table ").AddUnsafe(tableName))
// 	if err != nil {
// 		return fmt.Errorf("error executing drop table command: %w", err)
// 	}

// 	log.Println("table dropped successfully", tableName)
// 	return nil
// }
