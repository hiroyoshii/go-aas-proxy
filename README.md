<div align="center">
<h1>go-aas-proxy</h1>
<p>
Asset Administration Shell(AAS) proxy for RDBMS backend
</p>

[![Build](https://github.com/hiroyoshii/go-aas-proxy/actions/workflows/go_application.yaml/badge.svg)](https://github.com/hiroyoshii/go-aas-proxy/actions/workflows/go_application.yaml)
[![CodeQL](https://github.com/hiroyoshii/go-aas-proxy/actions/workflows/codeql.yml/badge.svg)](https://github.com/hiroyoshii/go-aas-proxy/actions/workflows/codeql.yml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/hiroyoshii/go-aas-proxy/badge)](https://securityscorecards.dev/viewer/?uri=github.com/hiroyoshii/go-aas-proxy)
[![License](https://img.shields.io/github/license/hiroyoshii/go-aas-proxy)](LICENSE)

</div>


## About

The go-aas-proxy is asset administration shell(aas) server which is implemented by go and proxy to RDBMS backend. **not ready for production**.

This implementation is inspired from "aas-proxy" in [Representing the Virtual: Using AAS to Expose DigitalAssets](https://ceur-ws.org/Vol-3291/paper5.pdf) paper.
![aas-proxy](./assets/aas-proxy.png)

"aas-proxy" means that submodels are configured by referencing the RDB of other applications.
Therefore, only the AAS and the AAS-Submodel relationship can be created, updated, or deleted, while submodels are read-only.

## Features
- Open API Endpoints compatible with [Basyx API](https://app.swaggerhub.com/apis/BaSyx/basyx_asset_administration_shell_http_rest_api/v1). But only supports following endpoints:
  - :white_check_mark: /shells: GET
  - :white_check_mark: /shells/{aasId}: GET/PUT/DELETE
  - :white_check_mark: /shells/{aasId}/aas/submodels: GET
  - :white_check_mark: /shells/{aasId}/aas/submodels/{submodelIdShort}: GET
  - :warning:	 /shells/{aasId}/aas/submodels/{submodelIdShort}: PUT(only relation to aas)
- Support multiple RDBMS locations and types (Postgres, MySQL)

Unsuported:
- :x: create/update/delete submodels and submodel elements
- :x: Invocation Endpoints

## Implementation Overview
![architecture](./assets/architecture.drawio.png)
- configurablity for aas
  - query sql and stored tables (default: aas contents are stored as json type in RDBMS)
- configurablity for submodel
  - reference databases per semanticID of submodel
  - query for databases
  - response json content for submodels

## How to run
* docker run
```
docker run hiroyoshii/go-aas-proxy:latest
```
* demo (docker-compose and demo deta)
```
./e2e/scenario_demo.sh
```

## License
Apache License 2.0, see [LICENSE](./LICENSE).
