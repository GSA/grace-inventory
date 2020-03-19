parts = $(subst -, ,$(CIRCLE_USERNAME))
environment := $(shell echo "$(word 2,$(parts))" | tr '[:lower:]' '[:upper:]')
environments := PRODUCTION STAGING TEST

ifeq ($(filter $(environment),$(environments)),)
	export environment = DEVELOPMENT
endif

export appenv := $(shell echo "$(environment)" | tr '[:upper:]' '[:lower:]')
export TF_VAR_appenv := $(appenv)

.PHONY: precommit test deploy check lint_handler test_handler build_handler release_handler integration_test plan_terraform validate_terraform init_terraform apply_terraform apply_terraform_tests destroy_terraform_tests clean
test: test_handler plan_terraform

deploy: build_handler apply_terraform

check: precommit
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

lint_handler: precommit
	make -C handler lint

test_handler: precommit
	make -C handler test

build_handler: precommit
	make -C handler build

release_handler: precommit
	make -C handler release

integration_test: precommit
	make -C tests integration_test

plan_terraform: validate_terraform
	terraform plan
	make -C tests plan

validate_terraform: init_terraform
	terraform validate

init_terraform: check
	[[ -d release ]] || mkdir release
	[[ -e release/grace-inventory-lambda.zip ]] || touch release/grace-inventory-lambda.zip
	terraform init -backend-config="bucket=$(backend_bucket)" -backend-config="key=$(backend_key)"

apply_terraform: apply_terraform_tests

apply_terraform_tests:
	make -C tests apply

destroy_terraform_tests:
	make -C tests destroy

clean: precommit
	make -C handler clean

precommit:
ifneq ($(strip $(hooksPath)),.github/hooks)
	@git config --add core.hooksPath .github/hooks
endif
