package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"google.golang.org/grpc"
)

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler      = (*AstraDatasource)(nil)
	_ backend.CheckHealthHandler    = (*AstraDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*AstraDatasource)(nil)
)

func NewDatasource(ctx context.Context, s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
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
	if d.conn != nil {
		err := d.conn.Close()
		if err != nil {
			backend.Logger.Error("Error closing database connection: %s.", err)
		}
	}
}
