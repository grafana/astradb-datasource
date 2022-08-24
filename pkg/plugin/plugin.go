package plugin

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler      = (*AstraDatasource)(nil)
	_ backend.CheckHealthHandler    = (*AstraDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*AstraDatasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := models.LoadSettings(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	return &AstraDatasource{settings: settings}, nil
}

type AstraDatasource struct {
	settings models.Settings
	client   *client.StargateClient
	conn     *grpc.ClientConn
}

func (d *AstraDatasource) Dispose() {
	// Clean up datasource instance resources.
	d.conn.Close()
}

func (d *AstraDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	rawCql string
}

func (d *AstraDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// create data frame response.
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *AstraDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
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

func (d *AstraDatasource) connect() error {
	if d.conn != nil {
		return nil
	}

	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(d.settings.URI, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(d.settings.Token),
		),
	)
	if err != nil {
		return err
	}

	d.conn = conn
	return nil
}
