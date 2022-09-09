package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v2"
)

type QueryModel struct {
	RawCql    string
	Format    sqlds.FormatQueryOption
	ActualCql string
}

func LoadQueryModel(query backend.DataQuery) (*QueryModel, error) {
	qm := &QueryModel{}
	err := json.Unmarshal(query.JSON, qm)
	return qm, err
}
