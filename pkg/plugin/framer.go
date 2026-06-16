package plugin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/converters"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

type column struct {
	field     *data.Field
	converter data.FieldConverter
	kind      string
}

func Frame(res *pb.Response, qm models.QueryModel) (*data.Frame, error) {

	result := res.GetResultSet()
	if result == nil {
		return data.NewFrame("response", nil), nil
	}

	columns, fields := getColumns(result)

	frame := data.NewFrame("response", fields...)

	var notices []data.Notice

	for _, row := range result.Rows {
		var vals []any

		for i, col := range columns {
			raw := row.Values[i]
			val, err := getValue(col, raw)
			if err != nil {
				notices = append(notices, data.Notice{Severity: data.NoticeSeverityWarning, Text: err.Error()})
			}
			vals = append(vals, val)
		}

		frame.AppendRow(vals...)
	}

	frame.Meta = &data.FrameMeta{
		ExecutedQueryString:    qm.ActualCql,
		PreferredVisualization: data.VisTypeGraph,
		Notices:                notices,
	}

	if getFormat(qm.Format) == sqlutil.FormatOptionTable {
		frame.Meta.PreferredVisualization = data.VisTypeTable
		return frame, nil
	}

	if getFormat(qm.Format) == sqlutil.FormatOptionLogs {
		frame.Meta.PreferredVisualization = data.VisTypeLogs
		return frame, nil
	}

	if frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
		fillMode := &data.FillMissing{Mode: data.FillModeNull}

		if shouldSort(qm.RawCql) {
			sortByTime(frame)
		}

		frame, err := data.LongToWide(frame, fillMode)
		if err != nil {
			return nil, err
		}
		return frame, nil
	}

	return frame, nil
}

func shouldSort(cql string) bool {
	lower := strings.ToLower(cql)
	return !strings.Contains(lower, "order by")
}

func sortByTime(frame *data.Frame) {
	// if the "as time" alias is used, use that first
	for _, f := range frame.Fields {
		if f.Name == "time" {
			sorter := experimental.NewFrameSorter(frame, f)
			sort.Sort(sorter)
			return
		}
	}
	// "as time" not used, just choose the first time field
	for _, f := range frame.Fields {
		if f.Type() == data.FieldTypeNullableTime || f.Type() == data.FieldTypeTime {
			sorter := experimental.NewFrameSorter(frame, f)
			sort.Sort(sorter)
			return
		}
	}
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
		return column{
			field: data.NewField(col.Name, nil, []*string{}),
			converter: data.FieldConverter{
				Converter: func(v any) (any, error) {
					v1, err := translateType(v.(*pb.Value), col.Type)
					if err != nil {
						return nil, err
					}
					mapInterface, ok := v1.(map[any]any)
					if !ok {
						return nil, nil
					}
					mapString := make(map[string]any)
					for key, value := range mapInterface {
						strKey := fmt.Sprintf("%v", key)
						strValue := value
						mapString[strKey] = strValue
					}
					b, err := json.Marshal(mapString)
					if err != nil {
						return nil, err
					}
					str := string(b)
					return &str, err
				},
			},
		}
	case *pb.TypeSpec_List_, *pb.TypeSpec_Set_, *pb.TypeSpec_Tuple_:
		return column{
			field: data.NewField(col.Name, nil, []*string{}),
			converter: data.FieldConverter{
				Converter: func(v any) (any, error) {
					v1, err := translateType(v.(*pb.Value), col.Type)
					if err != nil {
						return nil, err
					}
					b, err := json.Marshal(v1)
					if err != nil {
						return nil, err
					}
					str := string(b)
					return &str, err
				},
			},
		}
	case *pb.TypeSpec_Udt_:
		// TODO
	}

	field := data.NewField(name, labels, []*string{})
	field.Config = config
	return column{
		field,
		converters.AnyToNullableString,
		kind,
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
	case pb.TypeSpec_INET:
		return newColumn[string](col.Name, config, data.FieldConverter{
			Converter: func(v any) (any, error) {
				if v1, ok := v.(*pb.Value); ok && v1 != nil {
					if val := v1.GetInet(); val != nil {
						s := fmt.Sprintf("%v", val.GetValue())
						return &s, nil
					}
				}
				return nil, nil
			},
		}, v.String())
	case pb.TypeSpec_TIMEUUID:
		return newColumn[string](col.Name, config, data.FieldConverter{
			Converter: func(v any) (any, error) {
				if v1, ok := v.(*pb.Value); ok {
					if u, err := client.ToTimeUUID(v1); err == nil && u != nil {
						val := u.String()
						return &val, nil
					}
				}
				return nil, nil
			},
		}, v.String())
	case pb.TypeSpec_ASCII:
		return newColumn[string](col.Name, config, data.FieldConverter{
			Converter: func(v any) (any, error) {
				if v1, ok := v.(*pb.Value); ok {
					if u, err := client.ToString(v1); err == nil {
						return &u, nil
					}
				}
				return nil, nil
			},
		}, v.String())
	case pb.TypeSpec_UUID:
		return newColumn[string](col.Name, config, data.FieldConverter{
			Converter: func(v any) (any, error) {
				if v1, ok := v.(*pb.Value); ok {
					if u, err := client.ToUUID(v1); err == nil && u != nil {
						val := u.String()
						return &val, nil
					}
				}
				return nil, nil
			},
		}, v.String())
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
	if col.converter.Converter != nil {
		return col.converter.Converter(raw)
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
