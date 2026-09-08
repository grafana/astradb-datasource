import { expect } from '@jest/globals';
import { formatSQL } from './formatSql';

describe('formatSQL', () => {
  it('should format plain SQL', () => {
    expect(formatSQL('select a from b')).toBe('select\n  a\nfrom\n  b');
  });

  it.each([
    ['a variable', 'select * from t where id = $host', '$host'],
    ['a braced variable', 'select * from ${table}', '${table}'],
    ['a braced variable with a format', 'select * from t where x = ${var:sqlstring}', '${var:sqlstring}'],
    ['a legacy variable', 'select * from [[legacy]] where a = 1', '[[legacy]]'],
    ['a macro', 'select * from t where $__timeFilter(ts)', '$__timeFilter(ts)'],
    ['a macro with nested parentheses', 'select * from t where $__timeFilter(coalesce(a, b))', '$__timeFilter(coalesce(a, b))'],
    ['a macro taking a variable', 'select $__timeGroup(ts, $__interval) from t', '$__timeGroup(ts, $__interval)'],
  ])('should preserve %s', (_label, query, token) => {
    expect(formatSQL(query)).toContain(token);
  });

  it('should preserve every token in a query that mixes them', () => {
    const formatted = formatSQL('select * from ${table} where id = $host and $__timeFilter(ts)');

    expect(formatted).toContain('${table}');
    expect(formatted).toContain('$host');
    expect(formatted).toContain('$__timeFilter(ts)');
  });
});
