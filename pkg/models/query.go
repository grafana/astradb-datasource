package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sqlds "github.com/grafana/sqlds/v2"
)

type QueryModel struct {
	QueryType string                   `json:"queryType"`
	RawCql    string                   `json:"rawCql"`
	Format    *sqlds.FormatQueryOption `json:"format"`
	Dataset   string                   `json:"dataset"`
	Table     string                   `json:"table"`
	ActualCql string
}

func LoadQuery(query backend.DataQuery) (*QueryModel, error) {
	qm := &QueryModel{}
	err := json.Unmarshal(query.JSON, qm)
	if qm.Format == nil {
		qm.Format = &defaultQueryFormat
	}
	return qm, err
}
