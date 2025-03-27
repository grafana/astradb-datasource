import {
  DataFrame,
  DataFrameView,
  DataQuery,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  ScopedVars,
  TimeRange,
} from '@grafana/data';
import { CompletionItemKind, DB, LanguageCompletionProvider, QueryFormat } from '@grafana/plugin-ui';
import {
  BackendDataSourceResponse,
  DataSourceWithBackend,
  FetchResponse,
  getBackendSrv,
  getTemplateSrv,
  toDataQueryResponse,
} from '@grafana/runtime';
import { uniqueId } from 'lodash';
import { lastValueFrom } from 'rxjs';
import { map } from 'rxjs/operators';
import { buildColumnQuery, buildTableQuery, showDatabases } from './components/metaQuery';
import { fetchColumns, fetchTables, getFunctions, getSqlCompletionProvider } from './components/sqlCompletionProvider';
import type { AstraQuery, AstraSettings } from './types';

export class DataSource extends DataSourceWithBackend<AstraQuery, AstraSettings> {
  annotations = {};
  db: DB;
  dataset?: string;

  completionProvider: LanguageCompletionProvider | undefined;

  constructor(private instanceSettings: DataSourceInstanceSettings<AstraSettings>) {
    super(instanceSettings);
    this.db = this.getDB();
  }

  applyTemplateVariables(query: AstraQuery, scopedVars: ScopedVars) {
    const rawSql = this.replace(query.rawSql || '', scopedVars) || '';
    return { ...query, rawSql };
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
    if (!query.rawSql) {
      return [];
    }
    const frame = await this.runQuery(query);
    if (frame.fields?.length === 0) {
      return [];
    }
    if (frame?.fields?.length === 1) {
      return frame?.fields[0]?.values.map((text) => ({ text, value: text }));
    }
    // convention - assume the first field is an id field
    const ids = frame?.fields[0]?.values;
    return frame?.fields[1]?.values.map((text, i) => ({ text, value: ids.get(i) }));
  }

  async fetchDatasets(): Promise<string[]> {
    this.dataset = undefined;
    const datasets = await this.runSql<string[]>(showDatabases(), { refId: 'datasets' });
    return datasets
      .map((t) => t[0])
      .filter((d) => !d.startsWith('system') && !d.startsWith('datastax_') && !d.startsWith('data_'));
  }

  async fetchTables(dataset?: string): Promise<string[]> {
    this.dataset = dataset;
    const tables = await this.runSql<string[]>(buildTableQuery(dataset), { refId: 'tables' });
    return tables.map((t) => t[0]);
  }

  async fetchFields(query: Partial<AstraQuery>) {
    let ds = this.dataset || query.dataset;
    if (!ds || !query.table) {
      return [];
    }
    const queryString = buildColumnQuery({ ...query, dataset: ds }, query.table!);
    const frame = await this.runSql<string[]>(queryString, { refId: 'fields' });
    const fields = frame.map((f) => ({ name: f[0], text: f[0], value: f[0], type: f[1], label: f[0] }));
    return fields;
  }

  async fetchMeta(path?: string) {
    // TODO: most of our sql datasources allow setting a default db in config settings
    const defaultDB = this.instanceSettings.jsonData.database;
    path = path?.trim();
    if (!path && defaultDB) {
      const tables = await this.fetchTables(defaultDB);
      return tables.map((t) => ({ name: t, completion: t, kind: CompletionItemKind.Class }));
    } else if (!path) {
      const datasets = await this.fetchDatasets();
      return datasets.map((d) => ({ name: d, completion: `${d}.`, kind: CompletionItemKind.Module }));
    } else {
      const parts = path.split('.').filter((s: string) => s);
      if (parts.length > 2) {
        return [];
      }
      if (parts.length === 1 && !defaultDB) {
        const tables = await this.fetchTables(parts[0]);
        return tables.map((t) => ({ name: t, completion: t, kind: CompletionItemKind.Class }));
      } else if (parts.length === 1 && defaultDB) {
        const fields = await this.fetchFields({ dataset: defaultDB, table: parts[0] });
        return fields.map((t) => ({ name: t.value, completion: t.value, kind: CompletionItemKind.Field }));
      } else if (parts.length === 2 && !defaultDB) {
        const fields = await this.fetchFields({ dataset: parts[0], table: parts[1] });
        return fields.map((t) => ({ name: t.value, completion: t.value, kind: CompletionItemKind.Field }));
      } else {
        return [];
      }
    }
  }

  async runSql<T extends string[]>(query: string, options?: RunSQLOptions) {
    const frame = await this.runMetaQuery({ rawSql: query, format: QueryFormat.Table, refId: options?.refId }, options);
    return new DataFrameView<T>(frame);
  }

  private runMetaQuery(request: Partial<AstraQuery>, options?: MetricFindQueryOptions): Promise<DataFrame> {
    // TODO - if we need this we need to update grafana to export it
    // https://github.com/grafana/grafana/blob/main/public/app/features/dashboard/index.ts
    // const range = getTimeSrv().timeRange();
    const refId = request.refId || 'meta';
    const queries: DataQuery[] = [{ ...request, datasource: request.datasource || this.getRef(), refId }];

    return lastValueFrom(
      getBackendSrv()
        .fetch<BackendDataSourceResponse>({
          url: '/api/ds/query',
          method: 'POST',
          data: {
            from: options?.range?.from.valueOf().toString(), // || range.from.valueOf().toString(),
            to: options?.range?.to.valueOf().toString(), // || range.to.valueOf().toString(),
            queries,
          },
          requestId: refId,
        })
        .pipe(
          map((res: FetchResponse<BackendDataSourceResponse>) => {
            const rsp = toDataQueryResponse(res, queries);
            return rsp.data[0];
          })
        )
    );
  }

  runQuery(request: Partial<AstraQuery>): Promise<DataFrame> {
    return new Promise((resolve) => {
      const req = {
        targets: [{ ...request, refId: request.refId || uniqueId() }],
      } as DataQueryRequest<AstraQuery>;
      this.query(req).subscribe((res: DataQueryResponse) => {
        resolve(res.data[0] || { fields: [] });
      });
    });
  }

  getSqlCompletionProvider(db: DB): LanguageCompletionProvider {
    if (this.completionProvider !== undefined) {
      return this.completionProvider;
    }

    const args = {
      getColumns: { current: (query: AstraQuery) => fetchColumns(db, query) },
      getTables: { current: (dataset?: string) => fetchTables(db, { dataset }) },
      fetchMeta: { current: (path?: string) => this.fetchMeta(path) },
      getFunctions: { current: () => getFunctions() },
    };
    this.completionProvider = getSqlCompletionProvider(args);
    return this.completionProvider;
  }

  getDB(): DB {
    if (this.db !== undefined) {
      return this.db;
    }
    return {
      datasets: () => this.fetchDatasets(),
      tables: (dataset?: string) => this.fetchTables(dataset),
      fields: (query: AstraQuery) => this.fetchFields(query),
      validateQuery: (query: AstraQuery, range?: TimeRange) =>
        Promise.resolve({ query, error: '', isError: false, isValid: true }),
      dsID: () => this.id,
      lookup: (path?: string) => this.fetchMeta(path),
      getSqlCompletionProvider: () => this.getSqlCompletionProvider(this.db),
      functions: async () => getFunctions(),
      labels: new Map([['dataset', 'Keyspace']]),
    };
  }
}

interface RunSQLOptions extends MetricFindQueryOptions {
  refId?: string;
}

interface MetricFindQueryOptions extends SearchFilterOptions {
  range?: TimeRange;
}

export interface SearchFilterOptions {
  searchFilter?: string;
}
