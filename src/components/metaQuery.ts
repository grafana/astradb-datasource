import type { AstraQuery } from 'types';

const DATASET_NAME = 'keyspace_name';
const DATASETS = 'keyspaces';
const SCHEMA = 'system_schema';

export function showDatabases() {
  return `select ${DATASET_NAME} from ${SCHEMA}.${DATASETS};`;
}

export function buildTableQuery(dataset?: string) {
  return `SELECT table_name FROM ${SCHEMA}.tables where ${DATASET_NAME} = '${dataset}';`;
}

export function buildColumnQuery(queryModel: Partial<AstraQuery>, table: string, type?: string, timeColumn?: string) {
  let query = `SELECT column_name, type FROM ${SCHEMA}.columns WHERE `;
  query += buildTableConstraint(queryModel, table);

  switch (type) {
    case 'time': {
      query += " AND type IN ('timestamp','time','bigint','int','double','float')";
      break;
    }
    case 'metric': {
      query += " AND type IN ('text','varchar')";
      break;
    }
    case 'value': {
      query += " AND type IN ('bigint','int','smallint','tinyint','double','decimal','float')";
      query += ' AND column_name <> ' + quoteIdentAsLiteral(timeColumn!);
      break;
    }
    case 'group': {
      query += " AND type IN ('text','varchar')";
      break;
    }
  }

  query += ' ORDER BY column_name';

  return query;
}

export function buildTableConstraint(queryModel: Partial<AstraQuery>, table: string) {
  let query = '';

  // check for schema qualified table
  if (table.includes('.')) {
    const parts = table.split('.');
    query = `${DATASET_NAME} = ` + quoteIdentAsLiteral(parts[0]);
    query += ' AND table_name = ' + quoteIdentAsLiteral(parts[1]);
    return query;
  } else {
    query = `${DATASET_NAME} = '${queryModel.dataset}' AND table_name = ` + quoteIdentAsLiteral(table);

    return query;
  }
}

export function quoteIdentAsLiteral(value: string) {
  return quoteLiteral(unquoteIdentifier(value));
}

function quoteLiteral(value: string) {
  return "'" + value.replace(/'/g, "''") + "'";
}

function unquoteIdentifier(value: string) {
  if (value[0] === '"' && value[value.length - 1] === '"') {
    return value.substring(1, value.length - 1).replace(/""/g, '"');
  } else {
    return value;
  }
}
