package plugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
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

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	err := d.connect()

	if err != nil {
		status = backend.HealthStatusError
		message = err.Error()
	}
	if err == nil {
		_, err := client.NewStargateClientWithConn(d.conn)
		if err != nil {
			status = backend.HealthStatusError
			message = err.Error()
		}
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil

}
