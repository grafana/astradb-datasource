package plugin

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/converters"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

type column struct {
	field     *data.Field
	converter data.FieldConverter
	kind      string
}

// FormatQueryOption defines how the user has chosen to represent the data
type FormatQueryOption uint32

const (
	// FormatOptionTimeSeries formats the query results as a timeseries using "WideToLong"
	FormatOptionTimeSeries FormatQueryOption = iota
	// FormatOptionTable formats the query results as a table using "LongToWide"
	FormatOptionTable
	// FormatOptionLogs sets the preferred visualization to logs
	FormatOptionLogs
)

func Frame(res *pb.Response, qm QueryModel) *data.Frame {

	result := res.GetResultSet()
	if result == nil {
		return data.NewFrame("response", nil)
	}

	columns, fields := getColumns(result)

	frame := data.NewFrame("response", fields...)

	for _, row := range result.Rows {

		var vals []interface{}
		var errors []error

		for i, col := range columns {
			raw := row.Values[i]
			val, err := getValue(col, raw)
			if err != nil {
				fmt.Println(err.Error())
				// nolint:staticcheck
				errors = append(errors)
			}
			vals = append(vals, val)
		}

		frame.AppendRow(vals...)
	}

	frame.Meta = &data.FrameMeta{
		ExecutedQueryString:    qm.RawCql,
		PreferredVisualization: data.VisTypeGraph,
	}

	if qm.Format == FormatOptionTable {
		frame.Meta.PreferredVisualization = data.VisTypeTable
		return frame
	}

	if qm.Format == FormatOptionLogs {
		frame.Meta.PreferredVisualization = data.VisTypeLogs
		return frame
	}

	if frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
		fillMode := &data.FillMissing{Mode: data.FillModePrevious}
		frame, err := data.LongToWide(frame, fillMode)
		if err != nil {
			return nil
		}
		return frame
	}

	return frame
}

func getColumns(result *pb.ResultSet) ([]column, []*data.Field) {
	var columns []column
	var fields []*data.Field

	for _, col := range result.Columns {
		col := NewColumn(col, col.Name, col.Name, "", nil)
		columns = append(columns, col)
		fields = append(fields, col.field)
	}

	return columns, fields
}

func NewColumn(col *pb.ColumnSpec, name string, alias string, kind string, labels data.Labels) column {
	config := &data.FieldConfig{
		DisplayName: col.Name,
	}

	switch col.Type.Spec.(type) {
	case *pb.TypeSpec_Basic_:
		return newBasicColumn(col, config)
	case *pb.TypeSpec_Map_:
		// TODO
	case *pb.TypeSpec_List_:
		// TODO
	case *pb.TypeSpec_Set_:
		// TODO
	case *pb.TypeSpec_Udt_:
		// TODO
	case *pb.TypeSpec_Tuple_:
		// TODO
	}

	field := data.NewField(name, labels, []*string{})
	field.Config = config
	return column{
		field,
		converters.AnyToNullableString,
		"",
	}
}

func newBasicColumn(col *pb.ColumnSpec, config *data.FieldConfig) column {
	switch v := col.Type.GetBasic(); v {
	case pb.TypeSpec_DATE:
		return newColumn[time.Time](col.Name, config, dateTimeConverter, v.String())
	case pb.TypeSpec_TEXT, pb.TypeSpec_VARCHAR:
		return newColumn[string](col.Name, config, converters.AnyToNullableString, v.String())
	case pb.TypeSpec_DECIMAL:
		return newColumn[float64](col.Name, config, DecimalToNullableFloat64, v.String())
	case pb.TypeSpec_INT:
		return newColumn[int64](col.Name, config, converters.Int64ToNullableInt64, v.String())
	case pb.TypeSpec_BOOLEAN:
		return newColumn[bool](col.Name, config, converters.BoolToNullableBool, v.String())
	case pb.TypeSpec_FLOAT:
		return newColumn[float64](col.Name, config, Float32ToNullableFloat64, v.String())
	case pb.TypeSpec_DOUBLE:
		return newColumn[float64](col.Name, config, converters.Float64ToNullableFloat64, v.String())
	case pb.TypeSpec_BIGINT:
		return newColumn[int64](col.Name, config, BigIntConverter, v.String())
	case pb.TypeSpec_SMALLINT, pb.TypeSpec_TINYINT, pb.TypeSpec_COUNTER:
		return newColumn[int64](col.Name, config, SmallIntConverter, v.String())
	case pb.TypeSpec_VARINT:
		return newColumn[uint64](col.Name, config, VarIntConverter, v.String())
	case pb.TypeSpec_BLOB:
		return newColumn[string](col.Name, config, converters.AnyToNullableString, v.String())
	case pb.TypeSpec_TIME:
		return newColumn[uint64](col.Name, config, TimeConverter, v.String())
	case pb.TypeSpec_TIMESTAMP:
		return newColumn[time.Time](col.Name, config, TimestampConverter, v.String())
	default:
		return newColumn[string](col.Name, config, converters.AnyToNullableString, v.String())
	}
}

func getValue(col column, raw *pb.Value) (any, error) {
	switch col.kind {
	case pb.TypeSpec_DATE.String():
		return col.converter.Converter(raw.GetDate())
	case pb.TypeSpec_TEXT.String(), pb.TypeSpec_VARCHAR.String():
		return col.converter.Converter(raw.GetString_())
	case pb.TypeSpec_DECIMAL.String():
		return col.converter.Converter(raw)
	case pb.TypeSpec_INT.String():
		return col.converter.Converter(raw.GetInt())
	case pb.TypeSpec_BIGINT.String(), pb.TypeSpec_SMALLINT.String(), pb.TypeSpec_TINYINT.String(), pb.TypeSpec_VARINT.String(), pb.TypeSpec_COUNTER.String():
		return col.converter.Converter(raw)
	case pb.TypeSpec_BOOLEAN.String():
		return col.converter.Converter(raw.GetBoolean())
	case pb.TypeSpec_FLOAT.String():
		return col.converter.Converter(raw.GetFloat())
	case pb.TypeSpec_DOUBLE.String():
		return col.converter.Converter(raw.GetDouble())
	case pb.TypeSpec_TIME.String():
		return col.converter.Converter(raw.GetTime())
	case pb.TypeSpec_TIMESTAMP.String():
		return col.converter.Converter(raw.GetInt())
	case pb.TypeSpec_BLOB.String():
		v, err := client.ToBlob(raw)
		if err != nil {
			return nil, err
		}
		return col.converter.Converter(string(v))
	}
	return nil, nil
}

type Converted interface {
	float64 | int64 | int32 | uint64 | bool | string | time.Time
}

func newColumn[V Converted](name string, config *data.FieldConfig, converter data.FieldConverter, kind string) column {
	field := data.NewField(name, nil, []*V{})
	field.Config = config
	return column{
		field,
		converter,
		kind,
	}
}
