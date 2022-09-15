import { DataSource } from 'datasource';
import { SQLQuery } from 'plugin-ui';
import React, { useState } from 'react';
import { CQLEditor } from './CQLEditor';

export type CqlVariableQueryEditorProps = {
  datasource: DataSource;
  onChange: (query: SQLQuery, definition: string) => void;
  query: SQLQuery;
};

export const VariableQueryEditor = (props: CqlVariableQueryEditorProps) => {
  const { datasource, onChange } = props;
  const [cql, setCql] = useState<string>(props.query.rawSql || '');
  const [query, setQuery] = useState(props.query);
  const handleChange = (query: SQLQuery) => {
    setCql(query.rawSql || '');
    setQuery(query);
  };
  const onRun = () => {
    onChange({ ...query, rawSql: cql }, `Query: ${cql}`);
  };
  return (
    <CQLEditor datasource={datasource} onChange={handleChange} onRunQuery={onRun} query={{ ...query, rawSql: cql }} />
  );
};
