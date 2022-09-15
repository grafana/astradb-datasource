package models

import (
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v2"
)

type QueryModel struct {
	RawCql    string `json:"rawSql"`
	Format    any
	ActualCql string
}

func LoadQueryModel(query backend.DataQuery) (*QueryModel, error) {
	qm := &QueryModel{}
	err := json.Unmarshal(query.JSON, qm)
	if qm.Format == nil {
		qm.Format = sqlds.FormatOptionTable
	}
	if strings.Contains(strings.ToLower(qm.RawCql), "as time") {
		qm.Format = sqlds.FormatOptionTimeSeries
	}
	if strings.Contains(strings.ToLower(qm.RawCql), "as log_time") {
		qm.Format = sqlds.FormatOptionLogs
	}
	return qm, err
}
