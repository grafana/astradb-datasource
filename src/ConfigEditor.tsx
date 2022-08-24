import React, { ChangeEvent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { AstraSettings, SecureSettings } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<AstraSettings> {}

export const ConfigEditor = (props: Props) => {

  const onUriChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = props;
    const jsonData = {
      ...options.jsonData,
      uri: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  const onTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        token: event.target.value,
      },
    });
  };

  const onResetToken = () => {
    const { onOptionsChange, options } = props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
    });
  };

    const { options } = props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as SecureSettings;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="URI"
            labelWidth={6}
            inputWidth={20}
            onChange={onUriChange}
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
              inputWidth={20}
              onReset={onResetToken}
              onChange={onTokenChange}
            />
          </div>
        </div>
      </div>
    );
}
