import { OperatorType } from '@grafana/plugin-ui';

export const OPERATORS = [
  { type: OperatorType.Comparison, id: 'LESS_THAN', operator: '<', description: 'Returns TRUE if X is less than Y.' },
  {
    type: OperatorType.Comparison,
    id: 'LESS_THAN_EQUAL',
    operator: '<=',
    description: 'Returns TRUE if X is less than or equal to Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'GREATER_THAN',
    operator: '>',
    description: 'Returns TRUE if X is greater than Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'GREATER_THAN_EQUAL',
    operator: '>=',
    description: 'Returns TRUE if X is greater than or equal to Y.',
  },
  { type: OperatorType.Comparison, id: 'EQUAL', operator: '=', description: 'Returns TRUE if X is equal to Y.' },
  {
    type: OperatorType.Comparison,
    id: 'NOT_EQUAL',
    operator: '!=',
    description: 'Returns TRUE if X is not equal to Y.',
  },
  {
    type: OperatorType.Comparison,
    id: 'NOT_EQUAL_ALT',
    operator: '<>',
    description: 'Returns TRUE if X is not equal to Y.',
  },
  { type: OperatorType.Logical, id: 'AND', operator: 'AND' },
  { type: OperatorType.Logical, id: 'OR', operator: 'OR' },
];

export const AGGREGATE_FNS = [
  {
    id: 'AVG',
    name: 'AVG',
    description: `AVG(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the average of non-NULL input values, or NaN if the input contains a NaN.`,
  },
  {
    id: 'COUNT',
    name: 'COUNT',
    description: `COUNT(*)  [OVER (...)]
  Returns the number of rows in the input.
  COUNT(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the number of rows with expression evaluated to any value other than NULL.
  `,
  },
  {
    id: 'MAX',
    name: 'MAX',
    description: `MAX(
    expression
  )
  [OVER (...)]
  Returns the maximum value of non-NULL expressions. Returns NULL if there are zero input rows or expression evaluates to NULL for all rows. Returns NaN if the input contains a NaN.
  `,
  },
  {
    id: 'MIN',
    name: 'MIN',
    description: `MIN(
    expression
  )
  [OVER (...)]
  Returns the minimum value of non-NULL expressions. Returns NULL if there are zero input rows or expression evaluates to NULL for all rows. Returns NaN if the input contains a NaN.
  `,
  },
  {
    id: 'SUM',
    name: 'SUM',
    description: `SUM(
    [DISTINCT]
    expression
  )
  [OVER (...)]
  Returns the sum of non-null values.
  If the expression is a floating point value, the sum is non-deterministic, which means you might receive a different result each time you use this function.
  `,
  },
];
