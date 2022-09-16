import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/AstraQueryEditor';
import type { AstraQuery, AstraSettings } from './types';
import { VariableQueryEditor } from './components/VariableQueryEditor';

export const plugin = new DataSourcePlugin<DataSource, AstraQuery, AstraSettings>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setVariableQueryEditor(VariableQueryEditor);
