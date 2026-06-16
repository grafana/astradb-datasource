package plugin

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

var dateTimeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeTime,
	Converter: func(v any) (any, error) {
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

// DecimalToNullableFloat64 returns an error if the input is not a float64.
var DecimalToNullableFloat64 = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(v any) (any, error) {
		// TODO - seems to be an issue with decimals in the stargate package
		// as a workaround they can convert to float in the cql:  CAST( x AS float)
		convert := func(val *pb.Value) (*float64, error) {
			dec, err := client.ToDecimal(val)
			if err != nil {
				return nil, err
			}

			str := dec.String()
			if float, err := strconv.ParseFloat(str, 64); err == nil {
				return &float, nil
			}
			return nil, errors.New("unable to convert decimal to float")
		}
		return toNullable(v.(*pb.Value), convert)
	},
}

// Float32ToNullableFloat64 converts float32 to float64
var Float32ToNullableFloat64 = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(v any) (any, error) {
		convert := func(val any) (*float64, error) {
			v, ok := v.(float32)
			if !ok {
				return nil, errors.New("failed converting to float64")
			}
			f64 := math.Round((float64(v) * 100)) / 100
			return &f64, nil
		}
		return anyToNullable(v, convert)
	},
}

// BigIntConverter converts bigInt to float64
var BigIntConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableInt64,
	Converter: func(v any) (any, error) {
		return convertBigInt(v.(*pb.Value))
	},
}

func convertBigInt(val *pb.Value) (*int64, error) {
	convert := func(val *pb.Value) (*int64, error) {
		v, err := client.ToBigInt(val)
		if err != nil {
			return nil, err
		}
		if v.IsInt64() {
			return toInt64(v.Int64()), nil
		}
		if v.IsUint64() {
			return toInt64(v.Uint64()), nil
		}
		intVal, err := strconv.Atoi(v.String())
		if err != nil {
			return nil, errors.New("could not convert BigInt")
		}
		return toInt64(intVal), nil
	}
	return toNullable(val, convert)
}

// SmallIntConverter converts smallInt to int64
var SmallIntConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableInt64,
	Converter: func(v any) (any, error) {
		return convertSmallInt(v.(*pb.Value))
	},
}

func convertSmallInt(val *pb.Value) (*int64, error) {
	convert := func(val *pb.Value) (*int64, error) {
		v, err := client.ToSmallInt(val)
		return &v, err
	}
	return toNullable(val, convert)
}

// VarIntConverter converts smallInt to int64
var VarIntConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableInt64,
	Converter: func(v any) (any, error) {
		return convertVarInt(v.(*pb.Value))
	},
}

func convertVarInt(val *pb.Value) (*uint64, error) {
	convert := func(val *pb.Value) (*uint64, error) {
		v, err := client.ToVarInt(val)
		return &v, err
	}
	return toNullable(val, convert)
}

// TimeConverter converts uint64 to nullable uint64
var TimeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableInt64,
	Converter: func(v any) (any, error) {
		convert := func(val any) (*uint64, error) {
			v := val.(uint64)
			return &v, nil
		}
		return anyToNullable(v, convert)
	},
}

// TimestampConverter converts uint64 to time
var TimestampConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableInt64,
	Converter: func(v any) (any, error) {
		return convertIntToTimestamp(v.(int64))
	},
}

func convertIntToTimestamp(val int64) (*time.Time, error) {
	convert := func(val int64) (time.Time, error) {
		return time.Unix(0, val*int64(time.Millisecond)).UTC(), nil
	}
	converted, err := convert(val)
	if err != nil {
		return nil, err
	}

	return &converted, nil
}

type Int interface {
	int | int64 | uint64
}

func toInt64[V Int](val V) *int64 {
	v := int64(val)
	return &v
}

type Number interface {
	*int | *int64 | *uint64 | *float64 | *time.Time
}

func toNullable[T Number](val *pb.Value, f func(val *pb.Value) (T, error)) (T, error) {
	if val == nil {
		return nil, nil
	}
	return f(val)
}

func anyToNullable[T Number](val any, f func(val any) (T, error)) (T, error) {
	if val == nil {
		return nil, nil
	}
	return f(val)
}

func translateType(value *pb.Value, spec *pb.TypeSpec) (interface{}, error) {
	switch spec.GetSpec().(type) {
	case *pb.TypeSpec_Basic_:
		return translateBasicType(value, spec)
	case *pb.TypeSpec_Map_:
		elements := make(map[interface{}]interface{})
		if c := get(value.GetCollection()); len(c.Elements) > 0 {
			for i := 0; i < len(c.Elements)-1; i += 2 {
				key, err := translateType(c.Elements[i], spec.GetMap().Key)
				if err != nil {
					return nil, err
				}
				mapVal, err := translateType(c.Elements[i+1], spec.GetMap().Value)
				if err != nil {
					return nil, err
				}
				elements[key] = mapVal
			}
		}
		return elements, nil
	case *pb.TypeSpec_List_:
		var elements []interface{}
		for i := range get(value.GetCollection()).Elements {
			element, err := translateType(value.GetCollection().Elements[i], spec.GetList().Element)
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}

		return elements, nil
	case *pb.TypeSpec_Set_:
		var elements []interface{}
		for _, element := range get(value.GetCollection()).Elements {
			element, err := translateType(element, spec.GetSet().Element)
			if err != nil {
				return nil, err
			}

			elements = append(elements, element)
		}

		return elements, nil
	case *pb.TypeSpec_Udt_:
		fields := map[string]interface{}{}
		for key, val := range get(value.GetUdt()).Fields {
			element, err := translateType(val, spec.GetUdt().Fields[key])
			if err != nil {
				return nil, err
			}

			fields[key] = element
		}

		return fields, nil
	case *pb.TypeSpec_Tuple_:
		var elements []interface{}
		numElements := len(get(spec.GetTuple()).Elements)
		for i := 0; i <= len(get(value.GetCollection()).Elements)-numElements; i++ {
			for j, typeSpec := range spec.GetTuple().Elements {
				element, err := translateType(value.GetCollection().Elements[i+j], typeSpec)
				if err != nil {
					return nil, err
				}
				elements = append(elements, element)
			}
		}

		return elements, nil
	}
	return nil, errors.New("unsupported type")
}

func translateBasicType(value *pb.Value, spec *pb.TypeSpec) (interface{}, error) {
	switch spec.GetBasic() {
	case pb.TypeSpec_CUSTOM:
		return client.ToBlob(value)
	case pb.TypeSpec_ASCII, pb.TypeSpec_TEXT, pb.TypeSpec_VARCHAR:
		return client.ToString(value)
	case pb.TypeSpec_BIGINT:
		return client.ToBigInt(value)
	case pb.TypeSpec_BLOB:
		return client.ToBlob(value)
	case pb.TypeSpec_BOOLEAN:
		return client.ToBoolean(value)
	case pb.TypeSpec_COUNTER:
		return client.ToInt(value)
	case pb.TypeSpec_DECIMAL:
		return client.ToDecimal(value)
	case pb.TypeSpec_DOUBLE:
		return client.ToDouble(value)
	case pb.TypeSpec_FLOAT:
		return client.ToFloat(value)
	case pb.TypeSpec_INT:
		return client.ToInt(value)
	case pb.TypeSpec_TIMESTAMP:
		return client.ToTimestamp(value)
	case pb.TypeSpec_UUID:
		return client.ToUUID(value)
	case pb.TypeSpec_VARINT:
		return client.ToVarInt(value)
	case pb.TypeSpec_TIMEUUID:
		return client.ToTimeUUID(value)
	case pb.TypeSpec_INET:
		return client.ToInet(value)
	case pb.TypeSpec_DATE:
		return client.ToDate(value)
	case pb.TypeSpec_TIME:
		return client.ToTime(value)
	case pb.TypeSpec_SMALLINT:
		return client.ToSmallInt(value)
	case pb.TypeSpec_TINYINT:
		return client.ToTinyInt(value)
	}

	return nil, errors.New("unsupported type")
}

func get[T any](val *T) *T {
	if val == nil {
		return new(T)
	}
	return val
}
