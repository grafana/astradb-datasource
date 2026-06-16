import { DataQueryResponse, FieldType } from '@grafana/data';
import * as grafanaRuntime from '@grafana/runtime';
import { expect } from '@jest/globals';
import { from } from 'rxjs';
import { DataSource } from './datasource';

const getDataSource = () => {
  return new DataSource({
    readOnly: true,
    jsonData: {
      baseUrl: 'https://test-datasource.com',
      authenticateWithToken: false,
    },
  } as any);
};

describe('DataSource', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should create data source', () => {
    expect(getDataSource()).toBeTruthy();
  });

  it('should return correct `metricFindQuery` result', () => {
    const ds = getDataSource();
    const queryResponse: DataQueryResponse = {
      data: [
        {
          fields: [
            {
              values: ['test_one', 'test_two', 'test_three'],
              type: FieldType.string,
              name: 'name',
              config: {},
            },
          ],
          length: 3,
        },
      ],
    };
    ds.query = () => from([queryResponse]);

    expect(
      ds.metricFindQuery({
        refId: 'A',
        rawSql: 'f',
      })
    ).resolves.toStrictEqual([
      { text: 'test_one', value: 'test_one' },
      { text: 'test_two', value: 'test_two' },
      { text: 'test_three', value: 'test_three' },
    ]);
  });

  it('should return correct `applyTemplateVariables` result', () => {
    jest.spyOn(grafanaRuntime, 'getTemplateSrv').mockImplementation(() => ({
      replace: () => 'result string after replace',
      getVariables: jest.fn(),
      updateTimeRange: jest.fn(),
      containsTemplate: jest.fn(),
    }));

    const ds = getDataSource();
    const query = {
      refId: 'A',
    };

    expect(ds.applyTemplateVariables(query, { var: { text: '', value: '' } })).toStrictEqual({
      ...query,
      rawSql: 'result string after replace',
    });
  });
});
