package plugin

import (
	"database/sql"
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	sqlds "github.com/grafana/sqlds/v2"
)

// AstraDriver implements the driver interface for macro interpolation
// sqlds provides default macros using sqlds.Interpolate
type AstraDriver struct {
}

func (d AstraDriver) Connect(backend.DataSourceInstanceSettings, json.RawMessage) (*sql.DB, error) {
	return nil, nil
}

func (d AstraDriver) Settings(backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{}
}

func (d AstraDriver) Macros() sqlds.Macros {
	return sqlds.Macros{}
}

func (d AstraDriver) Converters() []sqlutil.Converter {
	return nil
}
