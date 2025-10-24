package plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	sqlds "github.com/grafana/sqlds/v4"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

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

func (d *AstraDatasource) query(_ context.Context, _ backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	qm, err := models.LoadQueryModel(query)
	if err != nil {
		response.Error = err
		return response
	}

	if qm.RawCql == "" {
		notice := data.Notice{Severity: data.NoticeSeverityWarning, Text: "empty query"}
		frame := data.Frame{Name: "warn", Meta: &data.FrameMeta{Notices: []data.Notice{notice}}, RefID: query.RefID}
		response.Frames = data.Frames{&frame}
		return response
	}

	stargateClient, err := client.NewStargateClientWithConn(d.conn)
	if err != nil {
		// Looks like there's no actual error here so we'll leave this as a plugin error
		response.Error = err
		return response
	}

	queryToEvaluate := &sqlutil.Query{
		RawSQL:    qm.RawCql,
		TimeRange: query.TimeRange,
		Format:    sqlutil.FormatQueryOption(getFormat(qm.Format)),
	}
	qm.ActualCql, err = sqlds.Interpolate(BaseDriver{}, queryToEvaluate)
	if err != nil {
		response.Error = err
		response.ErrorSource = backend.ErrorSourceDownstream
		return response
	}

	selectQuery := &pb.Query{
		Cql: qm.ActualCql,
	}

	queryResponse, err := stargateClient.ExecuteQuery(selectQuery)
	if err != nil {
		response.Error = err
		response.ErrorSource = backend.ErrorSourceDownstream
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

	var conn *grpc.ClientConn
	var err error

	if d.settings.AuthKind == models.AuthTypeToken {
		conn, err = grpc.NewClient(d.settings.URI, grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithPerRPCCredentials(
				auth.NewStaticTokenProvider(d.settings.Token),
			),
		)
	} else {
		if d.settings.Secure {
			config = &tls.Config{}
			conn, err = grpc.NewClient(d.settings.GRPCEndpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)),
				grpc.WithPerRPCCredentials(
					auth.NewTableBasedTokenProvider(
						fmt.Sprintf("https://%s/v1/auth", d.settings.AuthEndpoint), d.settings.UserName, d.settings.Password,
					),
				),
			)
		} else {
			conn, err = grpc.NewClient(d.settings.GRPCEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(
					auth.NewTableBasedTokenProviderUnsafe(
						fmt.Sprintf("http://%s/v1/auth", d.settings.AuthEndpoint), d.settings.UserName, d.settings.Password,
					),
				),
			)
		}
	}

	if err != nil {
		// Marking these errors as downstream as they could be due to misconfigured credentials
		return errorsource.DownstreamError(err, false)
	}

	d.conn = conn
	return nil
}

func getFormat(v any) sqlutil.FormatQueryOption {
	if v == nil {
		return sqlutil.FormatOptionTable
	}

	if fmt, ok := v.(string); ok {
		if strings.EqualFold(fmt, "time_series") {
			return sqlutil.FormatOptionTimeSeries
		} else if strings.EqualFold(fmt, "logs") {
			return sqlutil.FormatOptionLogs
		}
	}
	// for backwards compatibility with old format
	if fmt, ok := v.(*sqlutil.FormatQueryOption); ok {
		return *fmt
	}
	if fmt, ok := v.(sqlutil.FormatQueryOption); ok {
		return fmt
	}
	return sqlutil.FormatOptionTable
}
