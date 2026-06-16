import React, { useCallback } from 'react';
import type { QueryEditorProps } from '@grafana/data';
import type { DataSource } from '../datasource';
import type { AstraQuery, AstraSettings } from '../types';
import { SqlQueryEditor, SqlDatasource } from '@grafana/plugin-ui';
// @ts-ignore
import AutoSizer from 'react-virtualized-auto-sizer';

type Props = QueryEditorProps<DataSource, AstraQuery, AstraSettings>;

export const QueryEditor = ({ query, datasource, onChange, onRunQuery, range }: Props) => {
  const processQuery = useCallback(
    (q: AstraQuery) => {
      if (isQueryValid(q) && onRunQuery) {
        onRunQuery();
      }
    },
    [onRunQuery]
  );

  const onQueryChange = (q: AstraQuery, process = false) => {
    onChange(q);
    if (process) {
      processQuery(q);
    }
  };

  return (
    <SqlQueryEditor
      query={query}
      datasource={datasource as unknown as SqlDatasource}
      onRunQuery={onRunQuery}
      onChange={onQueryChange}
      range={range}
    />
  );
};

const isQueryValid = (q: AstraQuery) => {
  return Boolean(q.rawSql);
};
