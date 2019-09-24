parts = $(subst -, ,$(CIRCLE_USERNAME))
environment := $(shell echo "$(word 2,$(parts))" | tr '[:lower:]' '[:upper:]')
environments := PRODUCTION STAGING TEST

ifeq ($(filter $(environment),$(environments)),)
	export environment = DEVELOPMENT
endif

export appenv := $(shell echo "$(environment)" | tr '[:upper:]' '[:lower:]')
export TF_VAR_appenv := $(appenv)

.PHONY: test deploy lint_handler test_handler build_handler release_handler plan_terraform apply_terraform clean
test: test_handler plan_terraform

deploy: build_handler apply_terraform

lint_handler:
	make -C handler lint

test_handler:
	make -C handler test

build_handler:
	make -C handler build

release_handler:
	make -C handler release

plan_terraform:
	make -C terraform plan

apply_terraform:
	make -C terraform apply

apply_terraform_tests:
	make -C terraform/tests apply

destroy_terraform_tests:
	make -C terraform/tests destroy

clean:
	make -C handler clean
