import { expect, test } from '@grafana/plugin-e2e';

const ASTRA_URI = '37cd49dc-2aa3-4b91-a5e6-443c74d84c0c-us-east1.apps.astra.datastax.com:443';

test.describe('Test ConfigEditor', () => {
  test('invalid credentials should return an error', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-astradb-datasource' });

    await page.getByPlaceholder('$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443').fill(ASTRA_URI);
    await expect(configPage.saveAndTest()).not.toBeOK();
  });
});
