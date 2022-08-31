import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface AstraQuery extends DataQuery {
  rawCql: string;
  format: Format;
}

export interface AstraSettings extends DataSourceJsonData {
  uri: string;
}

export interface SecureSettings {
  token?: string;
}

export interface AutoSizerProps {
  width: number;
  height: number;
}

export enum Format {
  TIMESERIES = 0,
  TABLE = 1,
  LOGS = 2,
}
