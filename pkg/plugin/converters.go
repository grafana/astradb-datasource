package plugin

import (
	"errors"
	"math"
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

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

// Float64ToNullableFloat64 returns an error if the input is not a float64.
var Float32ToNullableFloat64 = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(v interface{}) (interface{}, error) {
		var ptr *float64
		if v == nil {
			return ptr, nil
		}
		val, ok := v.(float32)
		if !ok {
			return ptr, errors.New("failed converting to")
		}
		f64 := math.Round((float64(val) * 100)) / 100
		ptr = &f64
		return ptr, nil
	},
}
