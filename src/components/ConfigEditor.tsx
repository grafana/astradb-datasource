import React, { ChangeEvent } from 'react';
import { InlineLabel, LegacyForms, RadioButtonGroup, Checkbox, InlineFormLabel } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOptionChecked } from '@grafana/data';
import type { AstraSettings, SecureSettings } from '../types';
import { css } from '@emotion/css';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<AstraSettings> {}

export enum Connection {
  TOKEN = 0,
  CREDENTIALS = 1,
}

const types = [
  { label: 'Token', value: Connection.TOKEN },
  { label: 'Credentials', value: Connection.CREDENTIALS },
];

export const styles = {
  check: css`
    margin-top: 5px;
  `,
};

export const ConfigEditor = (props: Props) => {
  const onSecureSettingChange = (setting: 'password' | 'token') => (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        [setting]: event.target.value,
      },
    });
  };

  const onSettingChange =
    (setting: 'uri' | 'authEndpoint' | 'authKind' | 'grpcEndpoint' | 'user') =>
    (event: ChangeEvent<HTMLInputElement>) => {
      const { onOptionsChange, options } = props;
      const jsonData = {
        ...options.jsonData,
        [setting]: event.target.value,
      };
      onOptionsChange({ ...options, jsonData });
    };

  const onReset = (setting: 'password' | 'token') => {
    const { onOptionsChange, options } = props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        [setting]: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        [setting]: '',
      },
    });
  };

  const setConnectionType = (type: number) => {
    onSettingChange('authKind')(asEvent(type));
  };

  const { options } = props;
  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as SecureSettings;
  const kind = jsonData.authKind || Connection.TOKEN;

  return (
    <div className="gf-form-group">
      <div className="gf-form">
        <InlineLabel width={12}>Authentication</InlineLabel>
        <RadioButtonGroup options={types} value={kind} onChange={(v) => setConnectionType(v!)} size={'md'} />
      </div>

      {kind === Connection.TOKEN && (
        <>
          <div className="gf-form">
            <FormField
              label="URI"
              labelWidth={6}
              inputWidth={30}
              onChange={(v) => onSettingChange('uri')(v)}
              value={jsonData.uri || ''}
              placeholder="$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443"
            />
          </div>

          <div className="gf-form-inline">
            <div className="gf-form">
              <SecretFormField
                isConfigured={(secureJsonFields && secureJsonFields.token) as boolean}
                value={secureJsonData.token || ''}
                label="Token"
                placeholder="AstraCS:xxxxx"
                labelWidth={6}
                inputWidth={30}
                onReset={() => onReset('token')}
                onChange={(v) => onSecureSettingChange('token')(v)}
              />
            </div>
          </div>
        </>
      )}

      {kind === Connection.CREDENTIALS && (
        <>
          <div className="gf-form">
            <FormField
              label="GRPC Endpoint"
              labelWidth={7}
              inputWidth={20}
              onChange={(v) => onSettingChange('grpcEndpoint')(v)}
              value={jsonData.grpcEndpoint || ''}
              placeholder="localhost:8090"
            />
          </div>
          <div className="gf-form">
            <FormField
              label="Auth Endpoint"
              labelWidth={7}
              inputWidth={20}
              onChange={(v) => onSettingChange('authEndpoint')(v)}
              value={jsonData.authEndpoint || ''}
              placeholder="localhost:8081"
            />
          </div>
          <div className="gf-form">
            <FormField
              label="User Name"
              labelWidth={7}
              inputWidth={20}
              onChange={(v) => onSettingChange('user')(v)}
              value={jsonData.user || ''}
              placeholder="localhost:8090"
            />
          </div>
          <div className="gf-form-inline">
            <div className="gf-form">
              <SecretFormField
                isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
                value={secureJsonData.password || ''}
                label="Password"
                placeholder="xxxxx"
                labelWidth={7}
                inputWidth={20}
                onReset={() => onReset('password')}
                onChange={(v) => onSecureSettingChange('password')(v)}
              />
            </div>
          </div>
          <div className="gf-form">
            <InlineFormLabel width={7}>Secure</InlineFormLabel>
            <div className={styles.check}>
              <Checkbox
                value={jsonData.secure}
                onChange={onUpdateDatasourceJsonDataOptionChecked(props, 'secure')}
                label=""
                size={10}
              />
            </div>
          </div>
        </>
      )}
    </div>
  );
};

const asEvent = (value: string | number): ChangeEvent<HTMLInputElement> => {
  return { target: { value } } as unknown as ChangeEvent<HTMLInputElement>;
};
