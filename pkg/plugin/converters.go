package plugin

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

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

// DecimalToNullableFloat64 returns an error if the input is not a float64.
var DecimalToNullableFloat64 = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(v interface{}) (interface{}, error) {
		// TODO - seems to be an issue with decimals in the stargate package
		// as a workaround they can convert to float in the cql:  CAST( x AS float)
		var ptr *float64
		if v == nil {
			return ptr, nil
		}
		val, ok := v.(*pb.Value)
		if !ok {
			return ptr, errors.New("unable to convert decimal to float")
		}

		dec, err := client.ToDecimal(val)
		if err != nil {
			return nil, err
		}

		str := dec.String()
		if float, err := strconv.ParseFloat(str, 64); err == nil {
			return &float, nil
		}
		return ptr, errors.New("unable to convert decimal to float")
	},
}

// Float32ToNullableFloat64 converts float32 to float64
var Float32ToNullableFloat64 = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(v interface{}) (interface{}, error) {
		var ptr *float64
		if v == nil {
			return ptr, nil
		}
		val, ok := v.(float32)
		if !ok {
			return ptr, errors.New("failed converting to float64")
		}
		f64 := math.Round((float64(val) * 100)) / 100
		ptr = &f64
		return ptr, nil
	},
}

// TODO - code to reference for converting these types

// func translateType(spec *pb.TypeSpec) (interface{}, error) {
// 	switch spec.GetSpec().(type) {
// 	case *pb.TypeSpec_Basic_:
// 		return translateBasicType(value, spec)
// 	case *pb.TypeSpec_Map_:
// 		elements := make(map[interface{}]interface{})

// 		for i := 0; i < len(value.GetCollection().Elements)-1; i += 2 {
// 			key, err := translateType(value.GetCollection().Elements[i], spec.GetMap().Key)
// 			if err != nil {
// 				return nil, err
// 			}
// 			mapVal, err := translateType(value.GetCollection().Elements[i+1], spec.GetMap().Value)
// 			if err != nil {
// 				return nil, err
// 			}
// 			elements[key] = mapVal
// 		}
// 		return elements, nil
// 	case *pb.TypeSpec_List_:
// 		var elements []interface{}

// 		for i := range value.GetCollection().Elements {
// 			element, err := translateType(value.GetCollection().Elements[i], spec.GetList().Element)
// 			if err != nil {
// 				return nil, err
// 			}
// 			elements = append(elements, element)
// 		}

// 		return elements, nil
// 	case *pb.TypeSpec_Set_:
// 		var elements []interface{}
// 		for _, element := range value.GetCollection().Elements {
// 			element, err := translateType(element, spec.GetSet().Element)
// 			if err != nil {
// 				return nil, err
// 			}

// 			elements = append(elements, element)
// 		}

// 		return elements, nil
// 	case *pb.TypeSpec_Udt_:
// 		fields := map[string]interface{}{}
// 		for key, val := range value.GetUdt().Fields {
// 			element, err := translateType(val, spec.GetUdt().Fields[key])
// 			if err != nil {
// 				return nil, err
// 			}

// 			fields[key] = element
// 		}

// 		return fields, nil
// 	case *pb.TypeSpec_Tuple_:
// 		var elements []interface{}
// 		numElements := len(spec.GetTuple().Elements)
// 		for i := 0; i <= len(value.GetCollection().Elements)-numElements; i++ {
// 			for j, typeSpec := range spec.GetTuple().Elements {
// 				element, err := translateType(value.GetCollection().Elements[i+j], typeSpec)
// 				if err != nil {
// 					return nil, err
// 				}

// 				elements = append(elements, element)
// 			}
// 		}

// 		return elements, nil
// 	}
// 	return nil, errors.New("unsupported type")
// }
