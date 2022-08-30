import { DataSource } from 'datasource';
import React, { useState } from 'react';
import { CQLEditor } from './CQLEditor';
import { AstraQuery } from 'types';

export type CqlVariableQueryEditorProps = {
  datasource: DataSource;
  onChange: (query: AstraQuery, definition: string) => void;
  query: AstraQuery;
};

export const VariableQueryEditor = (props: CqlVariableQueryEditorProps) => {
  const { datasource, onChange } = props;
  const [cql, setCql] = useState<string>(props.query.rawCql || '');
  const [query, setQuery] = useState(props.query);
  const handleChange = (query: AstraQuery) => {
    setCql(query.rawCql || '');
    setQuery(query);
  };
  const onRun = () => {
    onChange({ ...query, rawCql: cql }, `Query: ${cql}`);
  };
  return (
    <CQLEditor datasource={datasource} onChange={handleChange} onRunQuery={onRun} query={{ ...query, rawCql: cql }} />
  );
};
