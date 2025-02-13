import React, { useCallback, useEffect, useRef } from 'react';

import { LanguageCompletionProvider, SQLEditor } from '@grafana/plugin-ui';

import { formatSQL } from '../utils/formatSql';
import type { AstraQuery } from 'types';

import { css } from '@emotion/css';

type Props = {
  query: AstraQuery;
  onRunQuery: () => void;
  onChange: (value: AstraQuery, processQuery: boolean) => void;
  children?: (props: { formatQuery: () => void }) => React.ReactNode;
  width?: number;
  height?: number;
  completionProvider?: LanguageCompletionProvider;
};

export function CQLEditor({ children, onChange, onRunQuery, query, width, height, completionProvider }: Props) {
  // We need to pass query via ref to SQLEditor as onChange is executed via monacoEditor.onDidChangeModelContent callback, not onChange property
  const queryRef = useRef<AstraQuery>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

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

  const onSqlChange = (sql: string) => {
    if (sql.trim() !== '') {
      onChange({ ...query, rawSql: sql }, true);
      onRunQuery();
    }
  };

  const run = () => onSqlChange(query.rawSql || '');

  const onQueryChange = useCallback(
    (rawCql: string, processQuery: boolean) => {
      const newQuery = {
        ...queryRef.current,
        rawQuery: true,
        rawCql,
      };
      onChange(newQuery, processQuery);
    },
    [onChange]
  );

  return (
    <div className={styles.wrapper}>
      <a onClick={run} className={styles.run}>
        <i className="fa fa-play"></i>
      </a>
      <SQLEditor
        width={width}
        height={height}
        query={query.rawSql!}
        onChange={onQueryChange}
        language={{ id: 'sql', completionProvider, formatter: formatSQL }}
      >
        {children}
      </SQLEditor>
    </div>
  );
}
