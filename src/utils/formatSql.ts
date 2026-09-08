import { format } from 'sql-formatter';

// Macro calls are masked whole, arguments included, because masking only the name
// makes the formatter insert a space before the parenthesis and Grafana then stops
// recognizing the macro.
const GRAFANA_TEMPLATE_SYNTAX =
  /\$__[a-zA-Z0-9_]+\s*\([^()]*(?:\([^()]*\)[^()]*)*\)|\$__[a-zA-Z0-9_]+|\$\{[^}]*\}|\$[a-zA-Z0-9_]+|\[\[[^\]]*\]\]/g;

export function formatSQL(q: string) {
  const tokens: string[] = [];
  const masked = q.replace(GRAFANA_TEMPLATE_SYNTAX, (token) => {
    tokens.push(token);
    return `grafana_template_${tokens.length - 1}`;
  });

  return tokens.reduce(
    (sql, token, i) => sql.replaceAll(`grafana_template_${i}`, () => token),
    format(masked)
  );
}
