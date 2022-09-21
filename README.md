# Astra DB data source plugin for Grafana

The Astra DB plugin allows a direct connection to Astra DB to query and visualize data in Grafana.

This plugin provides a CQL editor to format and color code your CQL statements, along with auto complete when entering keyspaces, tables and fields.

## Beta

This plugin is currently in **Beta** development. Breaking changes could occur but are not expected.

## Installation

For detailed instructions on how to install the plugin on Grafana Cloud or
locally, please checkout the [Plugin installation docs](https://grafana.com/docs/grafana/latest/plugins/installation/).

## Configuration

Once the plugin is installed on your Grafana instance, follow [these
instructions](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/)
to add a new Astra DB data source, and enter configuration options.

### Connect to Astra DB with authentication token
| Key   | Description                                                                                                                                                                                                                                 |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| URI   | Provide an Astra DB URI in the following format: `$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443`                                                                                                                              |
| Token | Provide an Astra DB token in the following format: `AstraCS:xxxxx`. See [Manage Application Tokens](https://docs.datastax.com/en/astra-serverless/docs/manage/org/managing-org.html#_manage_application_tokens) for more on authentication. |

### Connect to Cassandra with credentials ( requires [Stargate](https://stargate.io/docs/latest/core-index.html) )

| Key   | Description                                                                                                                                                                                                                                 |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GRPC Endpoint | GRPC Endpoint to Stargate.  Example:  `localhost:8090` |
| Auth Endpoint | Authentication Endpoint to Stargate.  Example:  `localhost:8081` |
| User Name | Cassandra database user name |
| Password | Cassandra database password |
| Secure | Check if secure connection is required |

### Configure via provisioning file

It is possible to configure data sources using configuration files with
Grafana’s provisioning system. To read about how it works, including all the
settings that you can set for this data source, refer to [Provisioning Grafana
data sources](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

Here is a provisioning example for this data source using token authentication:

```yaml
apiVersion: 1
datasources:
  - name: AstraDB
    type: grafana-astradb-datasource
    jsonData:
      uri: $ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443
    secureJsonData:
      token: AstraCS:xxxxx
```

### Time series

Time series visualization options are selectable after adding a `timestamp`
field type to your query. This field will be used as the timestamp. You can
select time series visualizations using the visualization options. Grafana
interprets timestamp rows without explicit time zone as UTC.

#### Multi-line time series

To create multi-line time series, the query must return at least 3 fields in
the following order:

- field 1: `timestamp` field with an alias of `time`
- field 2: value to group by
- field 3+: the metric values

For example:

```sql
SELECT time_field AS time, metric_name, avg(metric_value) AS avg_metric_value
FROM keyspace.table
GROUP BY metric_name, time_field
ORDER BY time_field
```

### Macros

To allowing injection of date range filters, the query can contain macros.

Here is an example of a query with a macros that will use the dashboard time range:

```sql
SELECT timestampvalue as time, bigintvalue, textvalue FROM grafana.tempTable1
where timestampvalue $__timeFrom and timestampvalue $__timeTo Allow Filtering
```

The query is converted to:

```sql
SELECT timestampvalue as time, bigintvalue, textvalue FROM grafana.tempTable1
where timestampvalue  >= '2021-07-07T12:04:16Z' and timestampvalue  <= '2021-11-08T21:26:04Z' Allow Filtering
```

### Templates and variables

To add a new query variable, refer to [Add a query
variable](https://grafana.com/docs/grafana/latest/variables/variable-types/add-query-variable/).

After creating a variable, you can use it in your CQL queries by using
[Variable syntax](https://grafana.com/docs/grafana/latest/variables/syntax/).
For more information about variables, refer to [Templates and
variables](https://grafana.com/docs/grafana/latest/variables/).

## Learn more

- Add [Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/).
- Configure and use [Templates and variables](https://grafana.com/docs/grafana/latest/variables/).
- Add [Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/).
- Set up alerting; refer to [Alerts overview](https://grafana.com/docs/grafana/latest/alerting/).
