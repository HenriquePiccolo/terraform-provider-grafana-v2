GRAFANA_VERSION ?= latest

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

testacc-enterprise: TESTARGS+=-tags='enterprise'
testacc-enterprise: testacc

testacc-cloud: TESTARGS+=-tags='cloud'
testacc-cloud: testacc

testacc-docker:
	GRAFANA_VERSION=$(GRAFANA_VERSION) \
		docker-compose \
		-f ./docker-compose.yml \
		run --rm grafana-provider \
		make testacc

testacc-docker-tls:
	GRAFANA_VERSION=$(GRAFANA_VERSION) \
		docker-compose \
		-f ./docker-compose.yml \
		-f ./docker-compose.tls.yml \
		run --rm grafana-provider \
		make testacc

changelog:
	@test $${RELEASE_VERSION?Please set environment variable RELEASE_VERSION}
	@test $${CHANGELOG_GITHUB_TOKEN?Please set environment variable CHANGELOG_GITHUB_TOKEN}
	@docker run -it --rm \
		-v $$PWD:/usr/local/src/your-app \
		-e CHANGELOG_GITHUB_TOKEN=$$CHANGELOG_GITHUB_TOKEN \
		ferrarimarco/github-changelog-generator \
		--user HenriquePiccolo \
		--project terraform-provider-grafana-v2 \
		--future-release $$RELEASE_VERSION
	@git add CHANGELOG.md && git commit -m "Release $$RELEASE_VERSION"

release:
	@test $${RELEASE_VERSION?Please set environment variable RELEASE_VERSION}
	@git tag $$RELEASE_VERSION
	@git push origin $$RELEASE_VERSION

drone:
	drone jsonnet --stream --source .drone/drone.jsonnet --target .drone/drone.yml
	drone lint .drone/drone.yml
	drone sign --save HenriquePiccolo/terraform-provider-grafana-v2 .drone/drone.yml
