- [ ] architecture + module
- [x] from open api to grpc

- [ ] pipeline
  - [ ] security pipeline, dependa bot
  - [ ] grpc linter
  - [ ] grpc gen code
  - [ ] k6 io testing
- [ ] badge
  - [ ] security badge


curl  https://app.swaggerhub.com/apis/BaSyx/basyx_asset_administration_shell_http_rest_api/v1/swagger.yaml?resolved=true > basyx-openapi.yaml
go install -v github.com/google/gnostic@latest
go install -v github.com/googleapis/gnostic-grpc@v0.1.0
gnostic --grpc-out=./proto basyx-openapi.yaml
