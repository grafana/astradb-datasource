import { DataQuery, DataSourceJsonData } from '@grafana/data';

//#region Settings
export interface AstraSettings extends DataSourceJsonData {
  uri: string;
}
export interface SecureSettings {
  token?: string;
}
//#endregion

//#region Query
export interface AstraQuery extends DataQuery {
  rawCql: string;
  format: Format;
}
export enum Format {
  TIMESERIES = 0,
  TABLE = 1,
  LOGS = 2,
}
//#endregion

//#region Misc
export interface AutoSizerProps {
  width: number;
  height: number;
}
//#endregion
