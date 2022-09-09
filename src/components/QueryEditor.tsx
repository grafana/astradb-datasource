import React, { useCallback, useMemo } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { Select, InlineFormLabel } from '@grafana/ui';
import { DataSource } from '../datasource';
import { CQLEditor } from './CQLEditor';
// @ts-ignore
import AutoSizer from 'react-virtualized-auto-sizer';
import { DEFAULT_QUERY_FORMAT, Format } from '../constants';
import type { AstraQuery, AstraSettings, AutoSizerProps } from '../types';

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
    <>
      <div className="gf-form">
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
      </div>
      <div className="gf-form">
        <InlineFormLabel>Format</InlineFormLabel>
        <Select<Format>
          options={[
            { label: 'Table', value: Format.TABLE },
            { label: 'TimeSeries', value: Format.TIMESERIES },
            { label: 'Logs', value: Format.LOGS },
          ]}
          value={query.format === undefined ? DEFAULT_QUERY_FORMAT : query.format}
          onChange={(e) => {
            if (e && e.value !== undefined) {
              onChange({ ...query, format: e.value });
              onRunQuery();
            }
          }}
        ></Select>
      </div>
    </>
  );
};

const isQueryValid = (q: AstraQuery) => {
  return Boolean(q.rawCql);
};
