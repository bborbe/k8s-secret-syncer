# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## Unreleased

- chore: update Go to 1.27.0 and github.com/bborbe/errors to v1.5.20, github.com/bborbe/http to v1.26.24, github.com/bborbe/k8s to v1.14.12, github.com/bborbe/log to v1.6.24, github.com/bborbe/metrics to v0.5.15, github.com/bborbe/run to v1.9.37, github.com/bborbe/sentry to v1.9.27, github.com/bborbe/service to v1.10.9, github.com/bborbe/time to v1.27.10, github.com/onsi/ginkgo/v2 to v2.32.1, k8s.io/api to v0.36.4, k8s.io/apimachinery to v0.36.4, k8s.io/client-go to v0.36.4

## v0.1.5

- fix: Emit a `version` label on `build_info` by migrating to shared `github.com/bborbe/metrics` and dropping the private `pkg/libmetrics` copy. The local copy registered a bare unlabelled `Gauge`, and the fleet `BuildStale` rule selects `build_info{version!~"^v[0-9]+[.][0-9]+[.][0-9]+$"}` — PromQL treats an absent label as `""`, which does not match, so every label-less series stayed permanently in scope and re-fired 14 days after each rebuild. Tagged releases now self-exempt while untagged builds (`v0.1.4-3-gabc1234`) stay correctly covered.

## v0.1.4

- chore: Bump errcheck to v1.20.0 and golangci-lint to v2.13.1 for Go 1.27 support
## v0.1.3

- chore: Ignore the generated /vendor/ dir — `make build` vendors it via check-go-mod (matching the trading repo's Makefile.docker), `make precommit` removes it again via ensure, and it is never tracked; leaving it unignored made `git add .` liable to commit the whole vendor tree
- chore: Update Go from 1.26.5 to 1.26.6 (go.mod, Dockerfile, CI) and golang.org/x/mod from v0.37.0 to v0.40.0 to clear the vulncheck failures blocking CI

## v0.1.2

- fix: mask SentryDSN in the startup argument dump with display:"length"
- fix: register the missing /setloglevel/{level} and /gc admin endpoints
- docs: add a License section to the README

## v0.1.1

- chore: Refresh build to clear the critical BuildStale alert on quant — the v0.1.0 image was 23 days old; BUILD_DATE and BUILD_GIT_VERSION are baked in at image build time, so only a rebuilt image resets build_info
- chore: Reformat the go.mod exclude directive to block form, as applied by go-modtool fmt in make precommit

## v0.1.0

- Extract k8s-secret-syncer from the trading monorepo into a dedicated public repo (publish-only image)
- Decouple from trading/lib (kafka-free metrics vendored into pkg/libmetrics)
