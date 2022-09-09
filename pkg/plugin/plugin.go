package plugin

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	sqlds "github.com/grafana/sqlds/v2"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler      = (*AstraDatasource)(nil)
	_ backend.CheckHealthHandler    = (*AstraDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*AstraDatasource)(nil)
)

func NewDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := models.LoadSettings(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	return &AstraDatasource{settings: settings}, nil
}

type AstraDatasource struct {
	settings models.Settings
	conn     *grpc.ClientConn
}

func (d *AstraDatasource) Dispose() {
	// Clean up datasource instance resources.
	d.conn.Close()
}

func (d *AstraDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	err := d.connect()
	if err != nil {
		return nil, err
	}

	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *AstraDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	qm, err := models.LoadQueryModel(query)
	if err != nil {
		response.Error = json.Unmarshal(query.JSON, qm)
		return response
	}

	stargateClient, err := client.NewStargateClientWithConn(d.conn)
	if err != nil {
		response.Error = err
		return response
	}

	queryToEvaluate := &sqlds.Query{
		RawSQL:    qm.RawCql,
		TimeRange: query.TimeRange,
		Format:    sqlds.FormatQueryOption(qm.Format),
	}
	qm.ActualCql, err = sqlds.Interpolate(BaseDriver{}, queryToEvaluate)
	if err != nil {
		response.Error = err
		return response
	}

	selectQuery := &pb.Query{
		Cql: qm.ActualCql,
	}

	queryResponse, err := stargateClient.ExecuteQuery(selectQuery)
	if err != nil {
		response.Error = err
		return response
	}
	frame, err := Frame(queryResponse, *qm)
	if err != nil {
		response.Error = err
		return response
	}
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *AstraDatasource) connect() error {
	// grpc - connect and stay connected
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
