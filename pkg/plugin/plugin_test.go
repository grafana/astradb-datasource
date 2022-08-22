package plugin_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// This is where the tests for the datasource backend live.
func TestQueryData(t *testing.T) {
	ds := plugin.AstraDatasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestConnect(t *testing.T) {

	t.Skip() // integration test

	// Astra DB configuration
	const astra_uri = "37cd49dc-2aa3-4b91-a5e6-443c74d84c0c-us-east1.apps.astra.datastax.com:443"
	const bearer_token = "AstraCS:LjDqrEIZyDgduvSZgHUKyfMX:25dc87b1f592f18d93261a45b13cd6b79a6bc43b9b79f7557749352030b62ea1"

	// Create connection with authentication
	// For Astra DB:
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(bearer_token),
		),
	)

	assert.Nil(t, err)
	assert.NotNil(t, conn)

	stargateClient, err := client.NewStargateClientWithConn(conn)

	assert.Nil(t, err)
	assert.NotNil(t, stargateClient)

	// For  Astra DB: SELECT the data to read from the table
	selectQuery := &pb.Query{
		Cql: "SELECT * from grafana.cars;",
	}

	response, err := stargateClient.ExecuteQuery(selectQuery)
	assert.Nil(t, err)

	result := response.GetResultSet()

	var i, j int
	for i = 0; i < 2; i++ {
		valueToPrint := ""
		for j = 0; j < 2; j++ {
			value, err := client.ToString(result.Rows[i].Values[j])
			if err != nil {
				fmt.Printf("error getting value %v", err)
			}
			valueToPrint += " "
			valueToPrint += value
		}
		fmt.Printf("%v \n", valueToPrint)
	}

}
