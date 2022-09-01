package plugin_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint     string
	authEndpoint     string
	astraDbContainer testcontainers.Container
)

func init() {
	ctx := context.Background()
	astraDbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "stargateio/stargate-3_11:v1.0.40",
			Env: map[string]string{
				"CLUSTER_NAME":    "test",
				"CLUSTER_VERSION": "3.11",
				"DEVELOPER_MODE":  "true",
				"ENABLE_AUTH":     "true",
			},
			ExposedPorts: []string{"8090/tcp", "8081/tcp", "8084/tcp", "9042/tcp"},
			WaitingFor:   wait.ForHTTP("/checker/readiness").WithPort("8084/tcp").WithStartupTimeout(90 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("Failed to start Stargate container: %v", err)
	}

	grpcPort, err := nat.NewPort("tcp", "8090")
	if err != nil {
		log.Fatalf("Failed to get port: %v", err)
	}

	authPort, err := nat.NewPort("tcp", "8081")
	if err != nil {
		log.Fatalf("Failed to get port: %v", err)
	}

	grpcEndpoint, err = astraDbContainer.PortEndpoint(ctx, grpcPort, "")
	if err != nil {
		log.Fatalf("Failed to get endpoint: %v", err)
	}

	authEndpoint, err = astraDbContainer.PortEndpoint(ctx, authPort, "")
	if err != nil {
		log.Fatalf("Failed to get endpoint: %v", err)
	}

}

func createClient(t *testing.T) *client.StargateClient {
	conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewTableBasedTokenProviderUnsafe(
				fmt.Sprintf("http://%s/v1/auth", authEndpoint), "cassandra", "cassandra",
			),
		),
	)
	require.NoError(t, err)

	astraDbClient, err := client.NewStargateClientWithConn(conn)
	require.NoError(t, err)
	return astraDbClient
}

//TODO: TestFramer currently tests for Frame() method returns correct frame
// But various data type fields returns null instead of actual data. This needs to be fixed
// Also later different types can be added
func TestFramer(t *testing.T) {
	stargateClient := createClient(t)

	// create keyspace
	query := &pb.Query{
		Cql: "CREATE KEYSPACE IF NOT EXISTS grafana WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};",
	}
	response, err := stargateClient.ExecuteQuery(query)
	require.NoError(t, err)

	assert.Nil(t, response.GetResultSet())

	// add table to keyspace
	cql := `
   CREATE TABLE IF NOT EXISTS grafana.tempTable1 (
     id uuid PRIMARY KEY,
     asciivalue ascii,
	 textvalue text,
	 varcharvalue varchar,
	 blobvalue blob,
	 booleanvalue boolean,
	 decimalvalue decimal,
	 doublevalue double,
  	 floatvalue float,
	 inetvalue inet,
     bigintvalue bigint,
	 intvalue int,
     smallintvalue smallint,
	 varintvalue varint,
	 tinyintvalue tinyint,
	 timevalue time,
	 timestampvalue timestamp,
     datevalue date,
     timeuuidvalue timeuuid,
     mapvalue map<int,text>,
     listvalue list<text>,
     setvalue set<text>,
     tuplevalue tuple<int, text, float>
   );`
	query = &pb.Query{
		Cql: cql,
	}
	response, err = stargateClient.ExecuteQuery(query)
	require.NoError(t, err)

	assert.Nil(t, response.GetResultSet())

	// insert into table
	cql = `
	INSERT INTO grafana.tempTable1 (
		id, 
		asciivalue,
		textvalue,
		varcharvalue,
		blobvalue,
		booleanvalue,
		decimalvalue,
		doublevalue,
		floatvalue,
		inetvalue,
		bigintvalue,
		intvalue,
		smallintvalue,
		varintvalue,
		tinyintvalue,
		timevalue,
		timestampvalue,
		datevalue,
		timeuuidvalue,
		mapvalue,
		listvalue,
		setvalue,
		tuplevalue
	) VALUES (
		f066f76d-5e96-4b52-8d8a-0f51387df76b,
		'alpha', 
		'bravo',
		'charlie',
		textAsBlob('foo'),
		true,
		1.1,
        2.2,
		3.3,
		'127.0.0.1',
        1,
		2,
		3,
		4,
		5,
        '10:15:30.123456789',
        '2021-09-07T16:40:31.123Z',
        '2021-09-07',
		30821634-13ad-11eb-adc1-0242ac120002,
		{1: 'a', 2: 'b', 3: 'c'},
		['a', 'b', 'c'],
		{'a', 'b', 'c'},
		(3, 'bar', 2.1)
	);
	`
	query = &pb.Query{
		Cql: cql,
	}
	response, err = stargateClient.ExecuteQuery(query)
	require.NoError(t, err)

	assert.Nil(t, response.GetResultSet())

	// read from table
	query = &pb.Query{
		Cql: "SELECT * FROM grafana.tempTable1",
	}
	response, err = stargateClient.ExecuteQuery(query)
	require.NoError(t, err)

	frameResponse := plugin.Frame(response)
	require.NotNil(t, frameResponse)
	experimental.CheckGoldenJSONFrame(t, "testdata", "framerAllTypes", frameResponse, updateGoldenFile)
}
