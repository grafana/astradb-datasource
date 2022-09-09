package models

import (
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v2"
)

type QueryModel struct {
	RawCql    string
	Format    *sqlds.FormatQueryOption
	ActualCql string
}

func LoadQueryModel(query backend.DataQuery) (*QueryModel, error) {
	qm := &QueryModel{
		Format: &TableFormat,
	}
	err := json.Unmarshal(query.JSON, qm)
	if strings.Contains(strings.ToLower(qm.RawCql), "as time") {
		qm.Format = &TimeSeriesFormat
	}
	return qm, err
}
