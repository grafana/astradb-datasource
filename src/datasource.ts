import {
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  ScopedVars,
  vectorator,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { uniqueId } from 'lodash';
import { AstraQuery, AstraSettings } from './types';

export class DataSource extends DataSourceWithBackend<AstraQuery, AstraSettings> {
  constructor(instanceSettings: DataSourceInstanceSettings<AstraSettings>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: AstraQuery, scopedVars: ScopedVars) {
    const sql = this.replace(query.rawCql || '', scopedVars) || '';
    return { ...query, rawCql: sql };
  }

  replace(value?: string, scopedVars?: ScopedVars) {
    if (value !== undefined) {
      return getTemplateSrv().replace(value, scopedVars, this.format);
    }
    return value;
  }

  format(value: any) {
    if (Array.isArray(value)) {
      return `'${value.join("','")}'`;
    }
    return value;
  }

  async metricFindQuery(query: AstraQuery) {
    if (!query.rawCql) {
      return [];
    }
    const frame = await this.runQuery(query);
    if (frame.fields?.length === 0) {
      return [];
    }
    if (frame?.fields?.length === 1) {
      return vectorator(frame?.fields[0]?.values).map((text) => ({ text, value: text }));
    }
    // convention - assume the first field is an id field
    const ids = frame?.fields[0]?.values;
    return vectorator(frame?.fields[1]?.values).map((text, i) => ({ text, value: ids.get(i) }));
  }

  runQuery(request: Partial<AstraQuery>): Promise<DataFrame> {
    return new Promise((resolve) => {
      const req = {
        targets: [{ ...request, refId: uniqueId() }],
      } as DataQueryRequest<AstraQuery>;
      this.query(req).subscribe((res: DataQueryResponse) => {
        resolve(res.data[0] || { fields: [] });
      });
    });
  }
}
