import { DataSourceJsonData } from '@grafana/data';
import { SQLQuery } from 'plugin-ui';

//#region Settings
export interface AstraSettings extends DataSourceJsonData {
  uri: string;
  database?: string;
}

export interface SecureSettings {
  token?: string;
}
//#endregion

export interface AstraQuery extends SQLQuery {}
