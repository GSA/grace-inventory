environment := integration-testing
export appenv := $(shell echo "$(environment)" | tr '[:upper:]' '[:lower:]')
export TF_VAR_appenv := $(appenv)
undefine TF_VAR_master_account_id
undefine TF_VAR_master_role_name
export backend_key := grace_integration_testing_inventory_lambda.tfstate
terraform := /usr/bin/terraform
terraform_version := 0.12.7

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

test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v ./...

$(terraform):
    curl https://releases.hashicorp.com/terraform/$(terraform_version)/terraform_$(terraform_version)_linux_amd64.zip -o /tmp/terraform_$(terraform_version)_linux_amd64.zip
    unzip -d /usr/bin/ /tmp/terraform_$(terraform_version)_linux_amd64.zip

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
