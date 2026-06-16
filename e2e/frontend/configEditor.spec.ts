import { expect, test } from '@grafana/plugin-e2e';

const ASTRA_URI = 'test';

test.describe('Test ConfigEditor', () => {
  test('invalid credentials should return an error', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-astradb-datasource' });

    await page.getByPlaceholder('$ASTRA_CLUSTER_ID-$ASTRA_REGION.apps.astra.datastax.com:443').fill(ASTRA_URI);
    await expect(configPage.saveAndTest()).not.toBeOK();
  });
});
