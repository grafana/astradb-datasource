# CONTRIBUTING

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

## Running E2E tests locally

1. Setup the provisioning file locally. (`provisioning/datasources/astradb.yaml`)
2. Build the plugin frontend and backend
3. Start the grafana `docker-compose up`
4. Execute `yarn e2e` in a separate console.

## Cluster details

You can spin up your own cluster using [AstraDB SAAS version](https://astra.datastax.com/) or use the one from the [provisioning repo](https://github.com/grafana/plugin-provisioning/blob/main/provisioning/datasources/astradb.yaml)
