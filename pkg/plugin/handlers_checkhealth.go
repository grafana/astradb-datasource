package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

func (d *AstraDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	if d.settings.URI == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Invalid AstraDB URL",
		}, nil
	}

	if d.settings.Token == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Invalid AstraDB Token",
		}, nil
	}

	err := d.connect()
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	c, err := client.NewStargateClientWithConn(d.conn)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	_, err = c.ExecuteQuery(&pb.Query{
		Cql: "select keyspace_name from system_schema.keyspaces;",
	})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil

}
