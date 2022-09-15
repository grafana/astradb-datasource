import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/AstraQueryEditor';
import { AstraSettings } from './types';
import { VariableQueryEditor } from './components/VariableQueryEditor';
import { SQLQuery } from 'plugin-ui';

export const plugin = new DataSourcePlugin<DataSource, SQLQuery, AstraSettings>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setVariableQueryEditor(VariableQueryEditor);
