SHELL := /bin/bash
default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:  ## make testacc
	source .env && TF_ACC=1 go test ./xenserver/ -v  $(TESTARGS) -timeout 120m

doc:  ## make doc for terraform provider documentation
	go generate ./...

provider: go.mod  ## make provider
	if [ -z "$(GOBIN)" ]; then echo "GOBIN is not set" && exit 1; fi
	go mod tidy
	go install .
	ls -l $(GOBIN)/terraform-provider-xenserver

apply: .env  ## make apply
	cd examples/terraform-main && \
    terraform plan && \
    terraform apply -auto-approve

show_state: .env  ## make show_state resource=xenserver_vm.vm
	@cd examples/terraform-main && \
	if [ -z "$(resource)" ]; then echo "USAGE: make show_state resource=<>" && \
	echo "List available resources:" && echo "`terraform state list`" && exit 1; fi && \
	terraform state show $(resource)

import: .env  ## make import resource=xenserver_vm.vm id=vm-uuid
	@cd examples/terraform-main && \
	if [ -z "$(resource)" ] || [ -z "$(id)" ]; then echo "USAGE: make import resource=<> id=<>"; exit 1; fi && \
	terraform import $(resource) $(id)

destroy:
	cd examples/terraform-main && \
    terraform destroy -auto-approve