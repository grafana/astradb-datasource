import { DataQuery, DataSourceJsonData, SelectableValue, TimeRange } from '@grafana/data';
import { CompletionItemKind, LanguageCompletionProvider, OperatorType } from '@grafana/experimental';

//#region Settings
export interface AstraSettings extends DataSourceJsonData {
  uri: string;
  database?: string;
}
export interface SecureSettings {
  token?: string;
}
//#endregion

//#region Query
export interface AstraQuery extends DataQuery {
  dataset: string;
  database?: string;
  rawCql: string;
  format: Format;
  table?: string;
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

//#region Query Editor

// TODO - pending grafana release - import this later from grafana
export interface MetaDefinition {
  name: string;
  completion?: string;
  kind: CompletionItemKind;
}

export interface Aggregate {
  id: string;
  name: string;
  description?: string;
}

export const OPERATORS = [
  { type: OperatorType.Comparison, id: 'LESS_THAN', operator: '<', description: 'Returns TRUE if X is less than Y.' },
  {
    type: OperatorType.Comparison,
    id: 'LESS_THAN_EQUAL',
    operator: '<=',
    description: 'Returns TRUE if X is less than or equal to Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'GREATER_THAN',
    operator: '>',
    description: 'Returns TRUE if X is greater than Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'GREATER_THAN_EQUAL',
    operator: '>=',
    description: 'Returns TRUE if X is greater than or equal to Y.',
  },
  { type: OperatorType.Comparison, id: 'EQUAL', operator: '=', description: 'Returns TRUE if X is equal to Y.' },
  {
    type: OperatorType.Comparison,
    id: 'NOT_EQUAL',
    operator: '!=',
    description: 'Returns TRUE if X is not equal to Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'NOT_EQUAL_ALT',
    operator: '<>',
    description: 'Returns TRUE if X is not equal to Y.',
  },
  { type: OperatorType.Logical, id: 'AND', operator: 'AND' },
  { type: OperatorType.Logical, id: 'OR', operator: 'OR' },
];

export const AGGREGATE_FNS = [
  {
    id: 'AVG',
    name: 'AVG',
    description: `AVG(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the average of non-NULL input values, or NaN if the input contains a NaN.`,
  },
  {
    id: 'COUNT',
    name: 'COUNT',
    description: `COUNT(*)  [OVER (...)]
  Returns the number of rows in the input.
  COUNT(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the number of rows with expression evaluated to any value other than NULL.
  `,
  },
  {
    id: 'MAX',
    name: 'MAX',
    description: `MAX(
    expression
  )
  [OVER (...)]
  Returns the maximum value of non-NULL expressions. Returns NULL if there are zero input rows or expression evaluates to NULL for all rows. Returns NaN if the input contains a NaN.
  `,
  },
  {
    id: 'MIN',
    name: 'MIN',
    description: `MIN(
    expression
  )
  [OVER (...)]
  Returns the minimum value of non-NULL expressions. Returns NULL if there are zero input rows or expression evaluates to NULL for all rows. Returns NaN if the input contains a NaN.
  `,
  },
  {
    id: 'SUM',
    name: 'SUM',
    description: `SUM(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the sum of non-null values.
  If the expression is a floating point value, the sum is non-deterministic, which means you might receive a different result each time you use this function.
  `,
  },
];

export interface DB {
  init?: (datasourceId?: string) => Promise<boolean>;
  datasets: () => Promise<string[]>;
  tables: (dataset?: string) => Promise<string[]>;
  fields: (query: AstraQuery, order?: boolean) => Promise<SQLSelectableValue[]>;
  validateQuery: (query: AstraQuery, range?: TimeRange) => Promise<ValidationResults>;
  dsID: () => number;
  dispose?: (dsID?: string) => void;
  lookup: (path?: string) => Promise<Array<{ name: string; completion: string }>>;
  getSqlCompletionProvider: () => LanguageCompletionProvider;
  toRawSql?: (query: AstraQuery) => string;
  functions: () => Promise<Aggregate[]>;
}

// React Awesome Query builder field types.
// These are responsible for rendering the correct UI for the field.
export type RAQBFieldTypes = 'text' | 'number' | 'boolean' | 'datetime' | 'date' | 'time';

export interface SQLSelectableValue extends SelectableValue {
  type?: string;
  raqbFieldType?: RAQBFieldTypes;
}

export interface ValidationResults {
  query: AstraQuery;
  rawSql?: string;
  error: string;
  isError: boolean;
  isValid: boolean;
  statistics?: {
    TotalBytesProcessed: number;
  } | null;
}

//#endregion
