package plugin_test

import (
	"testing"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: TestFramer currently tests for Frame() method returns correct frame
// But various data type fields returns null instead of actual data. This needs to be fixed
// Also later different types can be added
func TestFramer(t *testing.T) {
	stargateClient := createClient(t)

	// create keyspace
	query := &pb.Query{
		Cql: "CREATE KEYSPACE IF NOT EXISTS grafana WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};",
	}
	response, err := stargateClient.ExecuteQuery(query)
	require.NoError(t, err)
	assert.Nil(t, response.GetResultSet())

	response, err = createTestTable(stargateClient)
	require.NoError(t, err)
	assert.Nil(t, response.GetResultSet())

	_, err = insertTestData(stargateClient)
	require.NoError(t, err)

	// read from table
	query = &pb.Query{
		Cql: "SELECT * FROM grafana.tempTable1",
	}
	response, err = stargateClient.ExecuteQuery(query)
	require.NoError(t, err)

	qm := models.QueryModel{RawCql: query.Cql, Format: &models.TimeSeriesFormat}
	frameResponse, err := plugin.Frame(response, qm)
	require.Nil(t, err)
	require.NotNil(t, frameResponse)
	experimental.CheckGoldenJSONFrame(t, "testdata", "framerAllTypes", frameResponse, updateGoldenFile)
}

func createTestTable(stargateClient *client.StargateClient) (*pb.Response, error) {
	// add table to keyspace
	cql := `
	CREATE TABLE IF NOT EXISTS grafana.tempTable1 (
		id uuid PRIMARY KEY,
		asciivalue ascii,
		textvalue text,
		varcharvalue varchar,
		blobvalue blob,
		booleanvalue boolean,
		decimalvalue decimal,
		doublevalue double,
		floatvalue float,
		inetvalue inet,
		bigintvalue bigint,
		intvalue int,
		smallintvalue smallint,
		varintvalue varint,
		tinyintvalue tinyint,
		timevalue time,
		timestampvalue timestamp,
		datevalue date,
		timeuuidvalue timeuuid,
		mapvalue map<int,text>,
		listvalue list<text>,
		setvalue set<text>,
		tuplevalue tuple<int, text, float>
	);`
	query := &pb.Query{
		Cql: cql,
	}
	return stargateClient.ExecuteQuery(query)
}

func insertTestData(stargateClient *client.StargateClient) (*pb.Response, error) {
	// insert into table
	cql := `
		BEGIN BATCH
		INSERT INTO grafana.tempTable1 (
			id, 
			asciivalue,
			textvalue,
			varcharvalue,
			blobvalue,
			booleanvalue,
			decimalvalue,
			doublevalue,
			floatvalue,
			inetvalue,
			bigintvalue,
			intvalue,
			smallintvalue,
			varintvalue,
			tinyintvalue,
			timevalue,
			timestampvalue,
			datevalue,
			timeuuidvalue,
			mapvalue,
			listvalue,
			setvalue,
			tuplevalue
		) VALUES (
			f066f76d-5e96-4b52-8d8a-0f51387df763,
			'alpha', 
			'bravo',
			'charlie',
			textAsBlob('foo'),
			true,
			1.1,
			2.2,
			3.3,
			'127.0.0.1',
			5,
			2,
			3,
			4,
			5,
			'10:15:30.123456789',
			'2021-11-07T16:40:31.123Z',
			'2021-09-07',
			30821634-13ad-11eb-adc1-0242ac120002,
			{1: 'a', 2: 'b', 3: 'c'},
			['a', 'b', 'c'],
			{'a', 'b', 'c'},
			(3, 'bar', 2.1)
		)
		INSERT INTO grafana.tempTable1 (
			id, 
			asciivalue,
			textvalue,
			varcharvalue,
			blobvalue,
			booleanvalue,
			decimalvalue,
			doublevalue,
			floatvalue,
			inetvalue,
			bigintvalue,
			intvalue,
			smallintvalue,
			varintvalue,
			tinyintvalue,
			timevalue,
			timestampvalue,
			datevalue,
			timeuuidvalue,
			mapvalue,
			listvalue,
			setvalue,
			tuplevalue
		) VALUES (
			cbb8f69a-2a3a-11ed-a261-0242ac120003,
			'alpha', 
			'brave',
			'charlie',
			textAsBlob('foo'),
			true,
			1.2,
			2.2,
			3.3,
			'127.0.0.1',
			44,
			2,
			3,
			4,
			5,
			'10:15:30.123456789',
			'2021-11-07T16:40:31.123Z',
			'2021-09-07',
			30821634-13ad-11eb-adc1-0242ac120003,
			{1: 'a', 2: 'b', 3: 'c'},
			['a', 'b', 'c'],
			{'a', 'b', 'c'},
			(3, 'bar', 2.1)
		)
		INSERT INTO grafana.tempTable1 (
			id, 
			asciivalue,
			textvalue,
			varcharvalue,
			blobvalue,
			booleanvalue,
			decimalvalue,
			doublevalue,
			floatvalue,
			inetvalue,
			bigintvalue,
			intvalue,
			smallintvalue,
			varintvalue,
			tinyintvalue,
			timevalue,
			timestampvalue,
			datevalue,
			timeuuidvalue,
			mapvalue,
			listvalue,
			setvalue,
			tuplevalue
		) VALUES (
			16133404-2a3f-11ed-a261-0242ac120004,
			'alpha', 
			'brother',
			'charlie',
			textAsBlob('foo'),
			true,
			1.3,
			2.2,
			3.3,
			'127.0.0.1',
			444,
			2,
			3,
			4,
			5,
			'10:15:30.123456789',
			'2021-11-07T16:40:31.123Z',
			'2021-09-07',
			30821634-13ad-11eb-adc1-0242ac120002,
			{1: 'a', 2: 'b', 3: 'c'},
			['a', 'b', 'c'],
			{'a', 'b', 'c'},
			(3, 'bar', 2.1)
		)
		APPLY BATCH;
		`
	query := &pb.Query{
		Cql: cql,
	}
	return stargateClient.ExecuteQuery(query)
}
