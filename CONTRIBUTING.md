# CONTRIBUTING

## Signed commits are required

> [!IMPORTANT]
> All commits must be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) (GPG, SSH, or S/MIME) to be merged into this repository. Pull requests with unsigned commits will need to be re-committed with signatures before they can be merged.

## PreRequisites

- NodeJS version 16+
- Yarn latest
- Go 18+
- Mage 1.13.0+

## Setting up local dev environment

1. Clone the repo `git clone https://github.com/grafana/astradb-datasource` and `cd astradb-datasource`
2. Create the provisioning folder
   1. Create `provisioning/config/license/license.jwt` file and place the enterprise license
   2. Create a provisioning file `provisioning/datasources/astradb.yaml` and copy the content from [provisioning repo](https://github.com/grafana/plugin-provisioning/blob/main/provisioning/datasources/astradb.yaml)
3. Install npm dependencies `yarn install`
4. Build the frontend `yarn build`
5. Build the backend `mage -v`
6. Start Grafana `docker-compose up`
7. Open `http://localhost:3000`

## Data Source Configuration Schema

`pkg/schema/dsconfig.json` is the **single source of truth** for the data source's
configuration surface — every field a user can set, where it is stored (`jsonData` or
`secureJsonData`), its type, validation rules and UI hints. It is consumed by provisioning
tooling, documentation and automation.

The schema format is defined and documented by [`grafana/dsconfig`](https://github.com/grafana/dsconfig/tree/main/dsconfig):

- [README](https://github.com/grafana/dsconfig/tree/main/dsconfig#readme) — concepts and a worked example for each field shape (root / jsonData / secret / array / virtual), plus current gaps and limitations.
- [`schema.md`](https://github.com/grafana/dsconfig/blob/main/dsconfig/schema.md) — full property reference.
- [`schema.json`](https://github.com/grafana/dsconfig/blob/main/dsconfig/schema.json) — the JSON Schema `dsconfig.json` validates against. It is pinned via the `$schema` key at the top of our file, so editors autocomplete from it; bump that URL when you bump `github.com/grafana/dsconfig/schema` in `go.mod`.

The rest of this section covers only what is specific to this plugin.

### Layout

| File in `pkg/schema/`                | Description                                                                                                  |
| ------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `dsconfig.json`                       | Source of truth — **edit this**                                                                               |
| `dsconfig_test.go`                    | Wires the schema into the shared conformance suite; also holds `SecureKeys`                                   |
| `*.gen.json`                          | Generated artifacts — **never hand-edit**; `yarn build` copies them into `dist/schema/` via `webpack.config.ts` |

Note that the conformance suite checks the schema directly against `models.Settings`
(`pkg/models/settings.go`) — the same struct the backend reads settings into at runtime via
`LoadSettings`. There is no separate schema-only model to keep in sync; a field that affects
real backend behaviour and a field declared in the schema are always the same struct field.

### Adding a new settings option

1. **Declare the field** in `pkg/schema/dsconfig.json` under `fields`, and add its `id` to
   the appropriate `groups[].fieldRefs` entry. Field ids in this plugin follow the
   `<target>.<key>` convention, e.g. `jsonData.httpMethod`.
2. **Add the matching Go field** to `Settings` in `pkg/models/settings.go` with a json tag
   equal to the schema `key`. This parity is enforced in both directions — a field in the
   schema but not the struct (or vice versa) fails the test suite. Secrets
   (`target: secureJsonData`) are the exception: they get no `json` tag (use `json:"-"` and
   load them via `mapstructure` in `LoadSettings`, as `Token`/`Password` do), but their key
   must be added to `SecureKeys` in `pkg/schema/dsconfig_test.go`.
3. **If the field only applies to one authentication mode** (`jsonData.authKind`), set
   `dependsOn`/`requiredWhen` to `"jsonData.authKind == 0"` or `"jsonData.authKind == 1"`
   instead of a blanket `required: true` — see the existing auth fields for the pattern.
4. **Regenerate the artifacts** and commit them with your change:

   ```bash
   go generate ./pkg/schema/...
   ```

5. **Verify**:

   ```bash
   go test ./pkg/schema/...
   ```

### When the conformance suite fails

Most failures are self-explanatory from the assertion message. The three you are most
likely to hit:

- `SchemaArtifactInSync` — a `.gen.json` file has drifted. Run `go generate ./pkg/schema/...` and commit the result.
- `JSONDataMatchesStruct` / `JSONDataTypesMatchStruct` — the schema and `models.Settings` disagree on keys or types. Update whichever side is behind.
- `SecureValuesMatchLoadSettings` — the schema's `secureJsonData` fields and `SecureKeys` disagree.

## Running E2E tests locally

1. Setup the provisioning file locally. (`provisioning/datasources/astradb.yaml`)
2. Build the plugin frontend and backend
3. Start the grafana `docker-compose up`
4. Execute `yarn e2e` in a separate console.

## Cluster details

You can spin up your own cluster using [AstraDB SAAS version](https://astra.datastax.com/) or use the one from the [provisioning repo](https://github.com/grafana/plugin-provisioning/blob/main/provisioning/datasources/astradb.yaml)
