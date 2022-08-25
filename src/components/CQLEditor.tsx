import React, { useCallback, useEffect, useRef } from 'react';

import { LanguageCompletionProvider, SQLEditor } from '@grafana/experimental';

import { AstraQuery } from '../types';
import { formatSQL } from '../utils/formatSql';

type Props = {
  query: AstraQuery;
  onChange: (value: AstraQuery, processQuery: boolean) => void;
  children?: (props: { formatQuery: () => void }) => React.ReactNode;
  width?: number;
  height?: number;
  completionProvider?: LanguageCompletionProvider;
};

export function CQLEditor({ children, onChange, query, width, height, completionProvider }: Props) {
  // We need to pass query via ref to SQLEditor as onChange is executed via monacoEditor.onDidChangeModelContent callback, not onChange property
  const queryRef = useRef<AstraQuery>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

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
    <SQLEditor
      width={width}
      height={height}
      query={query.rawCql!}
      onChange={onQueryChange}
      language={{ id: 'sql', completionProvider, formatter: formatSQL }}
    >
      {children}
    </SQLEditor>
  );
}