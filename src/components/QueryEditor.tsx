import React, { useCallback } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { AstraQuery, AstraSettings, AutoSizerProps, Format } from '../types';
import { CQLEditor } from './CQLEditor';
// @ts-ignore
import AutoSizer from 'react-virtualized-auto-sizer';
import { css } from '@emotion/css';

type Props = QueryEditorProps<DataSource, AstraQuery, AstraSettings>;

export const QueryEditor = ({ query, onChange, onRunQuery }: Props) => {
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

  const onSqlChange = (sql: string) => {
    if (sql.trim() !== '') {
      const format = sql.toLowerCase().includes('as time') ? Format.TIMESERIES : Format.TABLE;
      onChange({ ...query, rawCql: sql, format });
      onRunQuery();
    }
  };

  const run = () => onSqlChange(query.rawCql || '');

  const styles = {
    wrapper: css`
      position: relative;
    `,
    run: css`
      position: absolute;
      top: 2px;
      left: 6px;
      z-index: 100;
      color: green;
    `,
  };

  return (
    <div style={{ width: '100%', height: '300px' }} className={styles.wrapper}>
      <a onClick={run} className={styles.run}>
        <i className="fa fa-play"></i>
      </a>
      <AutoSizer defaultHeight="300px" defaultWidth="100%">
        {(props: AutoSizerProps) => (
          <CQLEditor
            query={query}
            onChange={onQueryChange}
            width={props.width}
            height={props.height}
          />
        )}
      </AutoSizer>
    </div>
  );
}

const isQueryValid = (q: AstraQuery) => {
  return Boolean(q.rawCql);
};