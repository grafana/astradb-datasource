package plugin

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	sqlds "github.com/grafana/sqlds/v4"
)

// BaseDriver implements the driver interface for macro interpolation
// sqlds provides default macros using sqlds.Interpolate
type BaseDriver struct {
}

func (d BaseDriver) Connect(context.Context, backend.DataSourceInstanceSettings, json.RawMessage) (*sql.DB, error) {
	return nil, nil
}

func (d BaseDriver) Settings(context.Context, backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{}
}

func (d BaseDriver) Macros() sqlds.Macros {
	return sqlds.Macros{}
}

func (d BaseDriver) Converters() []sqlutil.Converter {
	return nil
}
