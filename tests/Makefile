environment := integration-testing
export appenv := $(shell echo "$(environment)" | tr '[:upper:]' '[:lower:]')
export TF_VAR_appenv := $(appenv)
undefine TF_VAR_master_account_id
undefine TF_VAR_master_role_name
export backend_key := grace_integration_testing_inventory_lambda.tfstate
<<<<<<< HEAD
=======
GO := /usr/local/go/bin/go
GO_VERSION := 1.13.5
>>>>>>> 021139995f847afd4bd7c85f08c4d173a65e971c

.PHONY: check apply plan validate init destroy test
check:
ifeq ($(strip $(backend_bucket)),)
	@echo "backend_bucket must be provided"
	@exit 1
endif
ifeq ($(strip $(TF_VAR_appenv)),)
	@echo "TF_VAR_appenv must be provided"
	@exit 1
else
	@echo "appenv: $(TF_VAR_appenv)"
endif
ifeq ($(strip $(backend_key)),)
	@echo "backend_key must be provided"
	@exit 1
endif

destroy: init
	terraform destroy -auto-approve

<<<<<<< HEAD
test:
	go test -v ./...
=======
test: $(GO)
	./usr/local/go/bin/go test -v ./...
>>>>>>> 021139995f847afd4bd7c85f08c4d173a65e971c

apply: plan
	terraform apply -auto-approve

plan: validate
	terraform plan

validate: init
	terraform validate
	terrascan --location . --tests all

init: check
	[[ -d ../release ]] || mkdir ../release
	[[ -e ../release/grace-inventory-lambda.zip ]] || touch ../release/grace-inventory-lambda.zip
	terraform init -backend-config="bucket=$(backend_bucket)" -backend-config="key=$(backend_key)"
<<<<<<< HEAD
=======

glibc:
	curl -L -o /tmp/glibc-2.30-r0.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-2.30-r0.apk
	curl -L -o /tmp/glibc-bin-2.30-r0.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-bin-2.30-r0.apk
	curl -L -o /tmp/glibc-dev-2.30-r0.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-dev-2.30-r0.apk
	curl -L -o /tmp/glibc-i18n-2.30-r0.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-i18n-2.30-r0.apk
	apk add --allow-untrusted /tmp/glibc-2.30-r0.apk /tmp/glibc-bin-2.30-r0.apk /tmp/glibc-dev-2.30-r0.apk /tmp/glibc-i18n-2.30-r0.apk

$(GO): glibc
	echo $(PATH)
	curl -L -o /tmp/go.tgz https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz
	tar -C /usr/local -xzf /tmp/go.tgz
	chmod +x $(GO)
>>>>>>> 021139995f847afd4bd7c85f08c4d173a65e971c
