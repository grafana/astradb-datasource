package models

import "github.com/grafana/sqlds/v2"

var (
	TableFormat      sqlds.FormatQueryOption = sqlds.FormatOptionTable
	TimeSeriesFormat sqlds.FormatQueryOption = sqlds.FormatOptionTimeSeries
)
