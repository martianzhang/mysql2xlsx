
# colors compatible setting
CRED:=$(shell tput setaf 1 2>/dev/null)
CGREEN:=$(shell tput setaf 2 2>/dev/null)
CYELLOW:=$(shell tput setaf 3 2>/dev/null)
CEND:=$(shell tput sgr0 2>/dev/null)

# Add mysql version for testing `MYSQL_RELEASE=percona MYSQL_VERSION=5.7 make docker`
# MySQL 5.1 `MYSQL_RELEASE=vsamov/mysql-5.1.73 make docker`
# MYSQL_RELEASE: mysql, percona, mariadb ...
# MYSQL_VERSION: latest, 8.0, 5.7, 5.6, 5.5 ...
# use mysql:latest as default
MYSQL_RELEASE := $(or ${MYSQL_RELEASE}, ${MYSQL_RELEASE}, mysql)
MYSQL_VERSION := $(or ${MYSQL_VERSION}, ${MYSQL_VERSION}, latest)

build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@mkdir -p bin
	@ret=0 && for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		go build -trimpath -o bin/$${b} $$d || ret=$$? ; \
	done ; exit $$ret
	@echo "build Success!"

GO_VERSION_MIN=1.13
.PHONY: go_version_check
# Parse out the x.y or x.y.z version and output a single value x*10000+y*100+z (e.g., 1.9 is 10900)
# that allows the three components to be checked in a single comparison.
VER_TO_INT:=awk '{split(substr($$0, match ($$0, /[0-9\.]+/)), a, "."); print a[1]*10000+a[2]*100+a[3]}'
go_version_check:
	@echo "$(CGREEN)Go version check ...$(CEND)"
	@if test $(shell go version | $(VER_TO_INT) ) -lt \
	$(shell echo "$(GO_VERSION_MIN)" | $(VER_TO_INT)); \
	then printf "go version $(GO_VERSION_MIN)+ required, found: "; go version; exit 1; \
		else echo "go version check pass";      fi

.PHONY: release
release: build
	@echo "$(CGREEN)Cross platform building for release ...$(CEND)"
	@mkdir -p release
	@for GOOS in darwin linux windows; do \
		for GOARCH in amd64; do \
			for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
				b=$$(basename $${d}) ; \
				echo "Building $${b}.$${GOOS}-$${GOARCH} ..."; \
				GOOS=$${GOOS} GOARCH=$${GOARCH} go build -trimpath -v -o release/$${b}.$${GOOS}-$${GOARCH} $$d 2>/dev/null ; \
			done ; \
		done ;\
	done

# Code format
.PHONY: fmt
fmt: go_version_check
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: docker-mysql
docker-mysql:
	@echo "$(CGREEN)Build mysql test environment ...$(CEND)"
	@docker stop xlsx-mysql 2>/dev/null || true
	@docker wait xlsx-mysql 2>/dev/null >/dev/null || true
	@echo "docker run --name xlsx-mysql $(MYSQL_RELEASE):$(MYSQL_VERSION)"
	@docker run --name xlsx-mysql --rm -d \
	-e MYSQL_ROOT_PASSWORD=123456 \
	-e MYSQL_DATABASE=test \
	-p 3306:3306 \
	$(MYSQL_RELEASE):$(MYSQL_VERSION)

	@echo "waiting for test database initializing "
	@timeout=180; while [ $${timeout} -gt 0 ] ; do \
		if ! docker exec xlsx-mysql mysql --user=root --password=123456 --host "127.0.0.1" --silent -NBe "do 1" >/dev/null 2>&1 ; then \
		        timeout=`expr $$timeout - 1`; \
		        printf '.' ;  sleep 1 ; \
		else \
		        echo "." ; echo "mysql test environment is ready!" ; break ; \
		fi ; \
		if [ $$timeout = 0 ] ; then \
		        echo "." ; echo "$(CRED)docker xlsx-mysql start timeout(180 s)!$(CEND)" ; exit 1 ; \
		fi ; \
	done

.PHONY: docker-connect
docker-connect:
	@docker exec -it xlsx-mysql mysql --user=root --password=123456 --host "127.0.0.1" test

# Run golang test cases
.PHONY: test
test:
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	go test -timeout 10m -race ./...
	@echo "test Success!"

# Code Coverage
# colorful coverage numerical >=90% GREEN, <80% RED, Other YELLOW
.PHONY: cover
cover: test
	@echo "$(CGREEN)Run test cover check ...$(CEND)"
	@go test $(LDFLAGS) -coverpkg=./... -coverprofile=test/coverage.data ./... | column -t
	@go tool cover -html=test/coverage.data -o test/coverage.html
	@go tool cover -func=test/coverage.data -o test/coverage.txt
	@tail -n 1 test/coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "$(CRED)"$$0"%$(CEND)"} \
		else if ($$NF >= 90) \
			{print "$(CGREEN)"$$0"%$(CEND)"} \
		else \
			{print "$(CYELLOW)"$$0"%$(CEND)"}}'

.PHONY: test-cli
test-cli: build
	@./bin/mysql2xlsx -user root --password 123456 \
	-query 'select "中文", "english"' \
	-file test/test-cli.csv

.PHONY: clean
clean:
	git clean -x -f