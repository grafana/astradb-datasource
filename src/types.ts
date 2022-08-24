import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface AstraQuery extends DataQuery {
  rawCql: string;
}

export interface AstraSettings extends DataSourceJsonData {
  uri: string;
}

export interface SecureSettings {
  token?: string;
}
