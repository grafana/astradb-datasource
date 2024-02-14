import React, { useState } from 'react';
import type { AstraQuery } from 'types';
import { CQLEditor } from './CQLEditor';

export type CqlVariableQueryEditorProps = {
  onChange: (query: AstraQuery, definition: string) => void;
  query: AstraQuery;
};

export const VariableQueryEditor = (props: CqlVariableQueryEditorProps) => {
  const { onChange } = props;
  const [cql, setCql] = useState<string>(props.query.rawSql || '');
  const [query, setQuery] = useState(props.query);
  const handleChange = (query: AstraQuery) => {
    setCql(query.rawSql || '');
    setQuery(query);
  };
  const onRun = () => {
    onChange({ ...query, rawSql: cql }, `Query: ${cql}`);
  };
  return (
    <CQLEditor onChange={handleChange} onRunQuery={onRun} query={{ ...query, rawSql: cql }} />
  );
};
