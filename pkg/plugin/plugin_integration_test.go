package plugin_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcEndpoint string
	authEndpoint string
)

// free tier - TODO - env vars
const astra_uri = "$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443"
const token = "AstraCS:xxxxx"
const updateGoldenFile = false

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
	os.Exit(0)
}

func setup() {
	_, shouldRun := os.LookupEnv("RUN_ASTRA_INTEGRATION_TESTS")
	if !shouldRun {
		fmt.Print("integration tests disabled.  set env var ")
		os.Exit(0)
	}

	_, runLocal := os.LookupEnv("RUN_LOCAL")
	if runLocal {
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

	_, dockerRunning := os.LookupEnv("DOCKER_RUNNING")
	if dockerRunning {
		grpcEndpoint = "localhost:8090"
		authEndpoint = "localhost:8081"
	}
}

func teardown() {
}

func TestConnect(t *testing.T) {
	// Create connection with authentication
	// For Astra DB:
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.NewClient(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(token),
		),
	)

	assert.Nil(t, err)
	assert.NotNil(t, conn)

	stargateClient, err := client.NewStargateClientWithConn(conn)

	assert.Nil(t, err)
	assert.NotNil(t, stargateClient)

	// For  Astra DB: SELECT the data to read from the table
	selectQuery := &pb.Query{
		Cql: "SELECT CAST( acceleration AS float) as acceleration, cylinders, displacement, horsepower, modelyear,  mpg,  passedemissions, CAST( weight as float) as weight from grafana.cars;",
	}

	response, err := stargateClient.ExecuteQuery(selectQuery)
	assert.Nil(t, err)

	qm := models.QueryModel{RawCql: selectQuery.Cql, Format: &models.TimeSeriesFormat}
	frame, err := plugin.Frame(response, qm)
	res := &backend.DataResponse{Frames: data.Frames{frame}, Error: err}

	experimental.CheckGoldenJSONResponse(t, "testdata", "connection", res, updateGoldenFile)
}

func TestQueryWithInts(t *testing.T) {
	r := runQuery(t, "SELECT show_id, date_added, release_year from grafana.movies_and_tv2 limit 10;")
	experimental.CheckGoldenJSONResponse(t, "testdata", "movies", r, updateGoldenFile)
}

func TestQueryWithTimeSeries(t *testing.T) {
	client := createRemoteClient(t)

	// createTestTable(client)
	_, err := insertTestData(client)
	require.NoError(t, err)

	// read from table
	query := &pb.Query{
		Cql: "SELECT timestampvalue as time, bigintvalue, textvalue FROM grafana.tempTable1",
	}
	response, err := client.ExecuteQuery(query)
	require.NoError(t, err)

	qm := models.QueryModel{RawCql: query.Cql, Format: &models.TimeSeriesFormat}
	frameResponse, err := plugin.Frame(response, qm)
	require.Nil(t, err)
	require.NotNil(t, frameResponse)
	experimental.CheckGoldenJSONFrame(t, "testdata", "timeseries", frameResponse, updateGoldenFile)
}

func runQuery(t *testing.T, cql string) *backend.DataResponse {
	query := fmt.Sprintf(`{"rawCql": "%s;"}`, cql)
	params := fmt.Sprintf(`{ "uri": "%s" }`, astra_uri)
	secure := map[string]string{"token": token}
	settings := backend.DataSourceInstanceSettings{JSONData: []byte(params), DecryptedSecureJSONData: secure}
	ds, err := plugin.NewDatasource(context.Background(), settings)
	assert.Nil(t, err)
	if err != nil {
		return nil
	}
	req := &backend.QueryDataRequest{
		Queries: []backend.DataQuery{
			{
				RefID:     "A",
				QueryType: "cql",
				JSON:      []byte(query),
			},
		},
		PluginContext: backend.PluginContext{
			DataSourceInstanceSettings: &settings,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dataSource := ds.(*plugin.AstraDatasource)
	res, err := dataSource.QueryData(ctx, req)
	assert.Nil(t, err)

	r := res.Responses["A"]
	return &r
}

func createClient(t *testing.T) *client.StargateClient {
	conn, err := grpc.NewClient(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func createRemoteClient(t *testing.T) *client.StargateClient {
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.NewClient(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(token),
		),
	)

	require.NoError(t, err)

	stargateClient, err := client.NewStargateClientWithConn(conn)
	require.NoError(t, err)
	return stargateClient
}

func TestConnectDocker(t *testing.T) {
	grpcEndpoint = "localhost:8090"
	authEndpoint = "localhost:8081"

	conn, err := grpc.NewClient(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(
			auth.NewTableBasedTokenProviderUnsafe(
				fmt.Sprintf("http://%s/v1/auth", authEndpoint), "cassandra", "cassandra",
			),
		),
	)
	if err != nil {
		log.Fatalf("error dialing connection %v", err)
	}

	stargateClient, err := client.NewStargateClientWithConn(conn)
	if err != nil {
		log.Fatalf("error creating client %v", err)
	}
	require.NoError(t, err)
	require.NotNil(t, stargateClient)

	err = conn.Close()
	require.NoError(t, err)
}
