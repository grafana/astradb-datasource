package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

func (d *AstraDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	if d.settings.URI == "" {
		return newHealthResult(backend.HealthStatusError, "Invalid AstraDB URL")
	}

	if d.settings.Token == "" {
		return newHealthResult(backend.HealthStatusError, "Invalid AstraDB Token")
	}

	if err := d.connect(); err != nil {
		return newHealthResult(backend.HealthStatusError, err.Error())
	}

	c, err := client.NewStargateClientWithConn(d.conn)
	if err != nil {
		return newHealthResult(backend.HealthStatusError, err.Error())
	}

	if _, err = c.ExecuteQuery(&pb.Query{
		Cql: "select keyspace_name from system_schema.keyspaces;",
	}); err != nil {
		return newHealthResult(backend.HealthStatusError, err.Error())
	}

	return newHealthResult(backend.HealthStatusOk, "Data source is working")
}

func newHealthResult(status backend.HealthStatus, message string) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
