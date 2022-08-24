import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { AstraQuery, AstraSettings } from './types';

export class DataSource extends DataSourceWithBackend<AstraQuery, AstraSettings> {
  constructor(instanceSettings: DataSourceInstanceSettings<AstraSettings>) {
    super(instanceSettings);
  }
}
