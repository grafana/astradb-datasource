package models

import (
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

type QueryModel struct {
	RawCql    string `json:"rawSql"`
	Format    any
	ActualCql string `json:"-"`
}

func LoadQueryModel(query backend.DataQuery) (*QueryModel, error) {
	qm := &QueryModel{}
	err := json.Unmarshal(query.JSON, qm)
	if qm.Format == nil {
		qm.Format = sqlutil.FormatOptionTable
	}
	if strings.Contains(strings.ToLower(qm.RawCql), "as time") {
		qm.Format = sqlutil.FormatOptionTimeSeries
	}
	if strings.Contains(strings.ToLower(qm.RawCql), "as log_time") {
		qm.Format = sqlutil.FormatOptionLogs
	}
	return qm, err
}
