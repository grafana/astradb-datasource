// Settle the masking design empirically before writing the real fix.
// Candidate A masks just the token head, candidate B masks whole macro calls too.
// A passing candidate must never throw and must return the Grafana tokens byte-identical.
import { format } from 'sql-formatter';

const A = /\$__[a-zA-Z0-9_]+|\$\{[^}]*\}|\$[a-zA-Z0-9_]+|\[\[[^\]]*\]\]/g;
const B = /\$__[a-zA-Z0-9_]+\s*\([^()]*(?:\([^()]*\)[^()]*)*\)|\$__[a-zA-Z0-9_]+|\$\{[^}]*\}|\$[a-zA-Z0-9_]+|\[\[[^\]]*\]\]/g;

function make(pattern) {
  return (q) => {
    const tokens = [];
    const masked = q.replace(pattern, (m) => {
      tokens.push(m);
      return `gf_token_${tokens.length - 1}`;
    });
    let out = format(masked);
    tokens.forEach((t, i) => {
      out = out.replaceAll(`gf_token_${i}`, t);
    });
    return out;
  };
}

const cases = [
  'select a from b',
  'select * from t where id = $host',
  'select * from ${table}',
  'select * from t where $__timeFilter(ts)',
  'select * from ${table} where id = $host and $__timeFilter(ts)',
  'select * from t where $__timeFilter(coalesce(a, b))',
  'select * from ${db.table} where x = ${var:sqlstring}',
  'select * from [[legacy]] where a = 1',
  'select $__timeGroup(ts, $__interval) , count(*) from t group by 1',
];

for (const [name, pattern] of [['A', A], ['B', B]]) {
  const fn = make(pattern);
  let pass = 0;
  console.log(`\n  ===== candidate ${name} =====`);
  for (const q of cases) {
    try {
      const out = fn(q);
      const oneline = out.replace(/\s+/g, ' ').trim();
      // Every Grafana token in the input must survive verbatim in the output.
      const tokens = q.match(/\$__[a-zA-Z0-9_]+\s*\([^()]*(?:\([^()]*\)[^()]*)*\)|\$\{[^}]*\}|\$[a-zA-Z0-9_]+|\[\[[^\]]*\]\]/g) || [];
      const lost = tokens.filter((t) => !out.includes(t));
      if (lost.length) {
        console.log(`    MANGLED ${JSON.stringify(q)}`);
        console.log(`            lost ${JSON.stringify(lost)} -> got ${JSON.stringify(oneline)}`);
      } else {
        pass++;
        console.log(`    ok      ${JSON.stringify(oneline)}`);
      }
    } catch (e) {
      console.log(`    THROWS  ${JSON.stringify(q)} -> ${String(e.message).split('\n')[0]}`);
    }
  }
  console.log(`    ${pass}/${cases.length} clean`);
}
