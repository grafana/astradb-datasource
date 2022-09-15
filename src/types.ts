import { DataSourceJsonData } from '@grafana/data';

//#region Settings
export interface AstraSettings extends DataSourceJsonData {
  uri: string;
  database?: string;
}

export interface SecureSettings {
  token?: string;
}
//#endregion
