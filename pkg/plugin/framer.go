package plugin

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/converters"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

type column struct {
	field     *data.Field
	converter data.FieldConverter
	kind      string
}

func Frame(res *pb.Response) *data.Frame {

	result := res.GetResultSet()
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
				errors = append(errors)
			}
			vals = append(vals, val)
		}

		frame.AppendRow(vals...)
	}
	return frame
}

func getColumns(result *pb.ResultSet) ([]column, []*data.Field) {
	var columns []column
	var fields []*data.Field

	for _, col := range result.Columns {
		col := NewColumn(col, "", "", "", nil)
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
		field := data.NewField(col.Name, nil, []*time.Time{})
		field.Config = config
		return column{
			field,
			dateTimeConverter,
			v.String(),
		}
	case pb.TypeSpec_TEXT, pb.TypeSpec_VARCHAR:
		field := data.NewField(col.Name, nil, []*string{})
		field.Config = config
		return column{
			field,
			converters.AnyToNullableString,
			v.String(),
		}
	case pb.TypeSpec_DECIMAL:
		field := data.NewField(col.Name, nil, []*float64{})
		field.Config = config
		return column{
			field,
			DecimalToNullableFloat64,
			v.String(),
		}
	case pb.TypeSpec_INT:
		field := data.NewField(col.Name, nil, []*int64{})
		field.Config = config
		return column{
			field,
			converters.Int64ToNullableInt64,
			v.String(),
		}
	case pb.TypeSpec_BOOLEAN:
		field := data.NewField(col.Name, nil, []*bool{})
		field.Config = config
		return column{
			field,
			converters.BoolToNullableBool,
			v.String(),
		}
	case pb.TypeSpec_FLOAT:
		field := data.NewField(col.Name, nil, []*float64{})
		field.Config = config
		return column{
			field,
			Float32ToNullableFloat64,
			v.String(),
		}
	case pb.TypeSpec_DOUBLE:
		field := data.NewField(col.Name, nil, []*float64{})
		field.Config = config
		return column{
			field,
			converters.Float64ToNullableFloat64,
			v.String(),
		}
	// TODO
	// pb.TypeSpec_BIGINT
	// pb.TypeSpec_BLOB
	// pb.TypeSpec_COUNTER
	// pb.TypeSpec_SMALLINT
	// pb.TypeSpec_TIME
	// pb.TypeSpec_TIMESTAMP
	// pb.TypeSpec_TINYINT
	// pb.TypeSpec_VARINT

	default:
		field := data.NewField(col.Name, nil, []*string{})
		field.Config = config
		return column{
			field,
			converters.AnyToNullableString,
			v.String(),
		}
	}
}

var dateTimeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeTime,
	Converter: func(v interface{}) (interface{}, error) {
		fV, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf(`expected %s input but got type %T for value "%v"`, "string", v, v)
		}
		t, err := dateparse.ParseAny(fV)
		if err != nil {
			return nil, fmt.Errorf("error converting to a time / date value. error: '%s', value: '%s", err.Error(), fV)
		}
		return &t, nil
	},
}

func getValue(col column, raw *pb.Value) (interface{}, error) {
	switch col.kind {
	case pb.TypeSpec_DATE.String():
		return col.converter.Converter(raw.GetDate())
	case pb.TypeSpec_TEXT.String(), pb.TypeSpec_VARCHAR.String():
		return col.converter.Converter(raw.GetString_())
	case pb.TypeSpec_DECIMAL.String():
		return col.converter.Converter(raw)
	case pb.TypeSpec_INT.String():
		return col.converter.Converter(raw.GetInt())
	case pb.TypeSpec_BOOLEAN.String():
		return col.converter.Converter(raw.GetBoolean())
	case pb.TypeSpec_FLOAT.String():
		return col.converter.Converter(raw.GetFloat())
	case pb.TypeSpec_DOUBLE.String():
		return col.converter.Converter(raw.GetDouble())
	}
	return nil, nil
}
