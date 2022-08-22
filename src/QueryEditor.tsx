import React, { ChangeEvent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { AstraQuery, AstraSettings } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, AstraQuery, AstraSettings>;

export const QueryEditor = (props: Props) => {
  const onQueryChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = props;
    onChange({ ...query, rawCql: event.target.value });
  };

  return (
    <div className="gf-form">
      <FormField
        labelWidth={8}
        value={props.query.rawCql || ''}
        onChange={onQueryChange}
        label="Query Text"
        tooltip="Not used yet"
      />
    </div>
  );
}
