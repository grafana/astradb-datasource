import React, { useCallback, useMemo } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { AstraQuery, AstraSettings, AutoSizerProps } from '../types';
import { CQLEditor } from './CQLEditor';
// @ts-ignore
import AutoSizer from 'react-virtualized-auto-sizer';

type Props = QueryEditorProps<DataSource, AstraQuery, AstraSettings>;

export const QueryEditor = ({ query, datasource, onChange, onRunQuery }: Props) => {
  const processQuery = useCallback(
    (q: AstraQuery) => {
      if (isQueryValid(q) && onRunQuery) {
        onRunQuery();
      }
    },
    [onRunQuery]
  );

  const onQueryChange = (q: AstraQuery, process = true) => {
    onChange(q);
    if (process) {
      processQuery(q);
    }
  };

  const completionProvider = useMemo(() => datasource.getDB().getSqlCompletionProvider(), [datasource]);

  return (
    <div style={{ width: '100%', height: '300px' }}>
      <AutoSizer defaultHeight="300px" defaultWidth="100%">
        {(props: AutoSizerProps) => (
          <CQLEditor
            query={query}
            datasource={datasource}
            onRunQuery={onRunQuery}
            onChange={onQueryChange}
            width={props.width}
            height={props.height}
            completionProvider={completionProvider}
          />
        )}
      </AutoSizer>
    </div>
  );
};

const isQueryValid = (q: AstraQuery) => {
  return Boolean(q.rawCql);
};
