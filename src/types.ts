import { DataSourceJsonData } from '@grafana/data';
import { SQLQuery } from '@grafana/plugin-ui';

//#region Settings
export interface AstraSettings extends DataSourceJsonData {
  uri: string;
  database?: string;
  authKind: number;
  user: string;
  password: string;
  grpcEndpoint: string;
  authEndpoint: string;
  secure: boolean;
}

export interface SecureSettings {
  token?: string;
  password?: string;
}
//#endregion

export interface AstraQuery extends SQLQuery {}
