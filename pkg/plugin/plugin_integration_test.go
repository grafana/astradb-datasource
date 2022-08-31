package plugin_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// free tier - TODO - env vars
const astra_uri = "37cd49dc-2aa3-4b91-a5e6-443c74d84c0c-us-east1.apps.astra.datastax.com:443"
const token = "AstraCS:LjDqrEIZyDgduvSZgHUKyfMX:25dc87b1f592f18d93261a45b13cd6b79a6bc43b9b79f7557749352030b62ea1"
const updateGoldenFile = false

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	_, shouldRun := os.LookupEnv("RUN_ASTRA_INTEGRATION_TESTS")
	if !shouldRun {
		os.Exit(0)
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

	conn, err := grpc.Dial(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
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

	frame := plugin.Frame(response)

	res := &backend.DataResponse{Frames: data.Frames{frame}, Error: err}

	experimental.CheckGoldenJSONResponse(t, "testdata", "connection", res, updateGoldenFile)
}

func TestQueryWithInts(t *testing.T) {
	r := runQuery(t, "SELECT show_id, date_added, release_year from grafana.movies_and_tv2 limit 10;")
	experimental.CheckGoldenJSONResponse(t, "testdata", "movies", r, updateGoldenFile)
}

func TestQueryWithTime(t *testing.T) {
	r := runQuery(t, "SELECT * FROM grafana.covidtime limit 10;")
	experimental.CheckGoldenJSONResponse(t, "testdata", "covidtime2", r, updateGoldenFile)
}

func TestQueryWithTimestamp(t *testing.T) {
	r := runQuery(t, "SELECT * FROM grafana.covid19 limit 10;")
	experimental.CheckGoldenJSONResponse(t, "testdata", "covid19", r, updateGoldenFile)
}

func runQuery(t *testing.T, cql string) *backend.DataResponse {
	query := fmt.Sprintf(`{"rawCql": "%s;"}`, cql)
	params := fmt.Sprintf(`{ "uri": "%s" }`, astra_uri)
	secure := map[string]string{"token": token}
	settings := backend.DataSourceInstanceSettings{JSONData: []byte(params), DecryptedSecureJSONData: secure}
	ds, err := plugin.NewDatasource(settings)
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
