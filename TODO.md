# Design
- [x] architecture + module
- [x] from open api to grpc

# Pipeline
- [ ] api
  - [x] gnostic & push
  - [x] buf gen go code
  - [x] buf lint
  - [x] buf breaking check
  - [x] buf format
- [ ] build
  - [ ] go build 
- [ ] testing
  - [ ] unit testing
  - [ ] senario testing
    - [ ] https://github.com/zoncoen/scenarigo
  - [ ] perf testing by k6.io
    - [ ] https://k6.io/blog/getting-started-with-performance-testing-in-ci-cd-using-k6/
- [ ] security
  - [ ] image scan
  - [ ] fuzzing
  - [ ] sbom
- [ ] image
  - [ ] image push
- [ ] badge
  - [ ] security badge: https://github.com/ossf/scorecard

# golang
- [ ] go1.21: https://tip.golang.org/doc/go1.21
  - [ ] slog
- [ ] aas
  - [ ] aas-core: 

# images
- [ ] docker compose


go install github.com/zoncoen/scenarigo/cmd/scenarigo@v0.14.2
scenarigo.exe run -c .\e2e\senariogo.yaml