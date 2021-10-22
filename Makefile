export PAGARMEAPI_ENDPOINT=https://api.pagar.me/1/transactions
export PAGARMEAPI_APIKEY=ak_test_IEIryoHZv7kppA0KWjkIDGGs5mucVl

VERSION_FLAG := -ldflags "-X github.com/pagarme/marshals/labs/vicki/desafioGo/config.version=${CIRCLE_SHA1}"

TEST_PACKAGES := $(shell find . -type f -iname \*_test.go | cut -d/ -f2- | xargs dirname | sort -u | xargs printf './%s ')

DC := docker-compose
GO := go

run:
	$(DC) up -d --force-recreate

log:
	$(DC) logs app

clean:
	$(DC) kill
	$(DC) rm -f
	docker rmi -f desafiogo_app

build:
	$(GO) build $(VERSION_FLAG)

test-with-coverage:
	for d in $(TEST_PACKAGES); do go test -coverprofile $$d.out $$d; done

sonar:
	docker run -ti -v $(shell pwd):/usr/src pagarme/sonar-scanner -Dsonar.branch.name=${BRANCH}