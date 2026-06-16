package models

import (
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

var (
	TableFormat      sqlutil.FormatQueryOption = sqlutil.FormatOptionTable
	TimeSeriesFormat sqlutil.FormatQueryOption = sqlutil.FormatOptionTimeSeries
)
