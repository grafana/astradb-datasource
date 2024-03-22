import { expect, test } from '@grafana/plugin-e2e';

const ASTRA_URI = '37cd49dc-2aa3-4b91-a5e6-443c74d84c0c-us-east1.apps.astra.datastax.com:443';
const ASTRA_TOKEN = 'AstraCS:LjDqrEIZyDgduvSZgHUKyfMX:25dc87b1f592f18d93261a45b13cd6b79a6bc43b9b79f7557749352030b62ea1';

test.describe('Test ConfigEditor', () => {
  test('invalid credentials should return an error', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: 'astradb-datasource' });
    await page.getByPlaceholder('$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443').fill(ASTRA_URI);
    await page.getByPlaceholder('AstraCS:xxxxx').fill(ASTRA_TOKEN);
    await expect(configPage.saveAndTest()).not.toBeOK();
  });
});
