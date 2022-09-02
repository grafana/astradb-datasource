import { e2e } from '@grafana/e2e';
import { PROVISIONING_FILE_NAME, PLUGIN_NAME } from './utils';

describe('config', () => {
  it('incomplete settings should throw error', () => {
    e2e.flows.login();
    e2e()
      .readProvisions([PROVISIONING_FILE_NAME])
      .then(([provision]) => {
        e2e.flows.addDataSource({
          expectedAlertMessage: 'Invalid AstraDB URL',
          form: () => {},
          type: PLUGIN_NAME,
        });
      });
  });
});
