package plugin_test

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/auth"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// This is where the tests for the datasource backend live.
func TestQueryData(t *testing.T) {
	ds := plugin.AstraDatasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestConnect(t *testing.T) {

	t.Skip() // integration test

	// Astra DB configuration
	const astra_uri = "37cd49dc-2aa3-4b91-a5e6-443c74d84c0c-us-east1.apps.astra.datastax.com:443"
	const bearer_token = "AstraCS:LjDqrEIZyDgduvSZgHUKyfMX:25dc87b1f592f18d93261a45b13cd6b79a6bc43b9b79f7557749352030b62ea1"

	// Create connection with authentication
	// For Astra DB:
	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := grpc.Dial(astra_uri, grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			auth.NewStaticTokenProvider(bearer_token),
		),
	)

	assert.Nil(t, err)
	assert.NotNil(t, conn)

	stargateClient, err := client.NewStargateClientWithConn(conn)

	assert.Nil(t, err)
	assert.NotNil(t, stargateClient)

	// For  Astra DB: SELECT the data to read from the table
	selectQuery := &pb.Query{
		Cql: "SELECT CAST( acceleration AS float) as acceleration, cylinders, displacement, horsepower, modelyear,  mpg,  passedemissions, CAST( weight as float) as weight from grafana.cars;",
	}

	response, err := stargateClient.ExecuteQuery(selectQuery)
	assert.Nil(t, err)

	frame := plugin.Frame(response)

	res := &backend.DataResponse{Frames: data.Frames{frame}, Error: err}

	err = experimental.CheckGoldenDataResponse("../testdata/basic.txt", res, true)
	assert.Nil(t, err)
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
