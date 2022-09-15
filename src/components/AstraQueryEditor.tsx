import React, { useCallback } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { AstraSettings } from '../types';
// import { CQLEditor } from './CQLEditor';
// @ts-ignore
import { SqlQueryEditor, SqlDatasource, SQLQuery, SQLOptions } from 'plugin-ui';
// @ts-ignore
import AutoSizer from 'react-virtualized-auto-sizer';

type Props = QueryEditorProps<DataSource, SQLQuery, AstraSettings>;

export const QueryEditor = ({ query, datasource, onChange, onRunQuery, range }: Props) => {
  const processQuery = useCallback(
    (q: SQLQuery) => {
      if (isQueryValid(q) && onRunQuery) {
        onRunQuery();
      }
    },
    [onRunQuery]
  );

  const onQueryChange = (q: SQLQuery, process = false) => {
    onChange(q);
    if (process) {
      processQuery(q);
    }
  };

  // const onRun = () => {
  //   onRunQuery();
  // }
  // const completionProvider = useMemo(() => datasource.getDB().getSqlCompletionProvider(), [datasource]);

  // type Props = QueryEditorProps<SqlDatasource, SQLQuery, SQLOptions>;

  return (
    // <div style={{ width: '100%', height: '300px' }}>
    //   <AutoSizer defaultHeight="300px" defaultWidth="100%">
    //     {(props: AutoSizerProps) => (
    //       <CQLEditor
    //         query={query}
    //         datasource={datasource}
    //         onRunQuery={onRunQuery}
    //         onChange={onQueryChange}
    //         width={props.width}
    //         height={props.height}
    //         completionProvider={completionProvider}
    //       />
    //     )}
    //   </AutoSizer>
    // </div>
    <SqlQueryEditor
      query={query}
      datasource={datasource as unknown as SqlDatasource}
      onRunQuery={onRunQuery}
      onChange={onQueryChange}
      range={range}
    />
  );
};

const isQueryValid = (q: SQLQuery) => {
  return Boolean(q.rawSql);
};
