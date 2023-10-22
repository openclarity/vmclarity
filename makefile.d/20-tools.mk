####
## Licensei CLI
####

LICENSEI_BIN := $(BIN_DIR)/licensei
LICENSEI_VERSION = 0.9.0

bin/licensei: bin/licensei-$(LICENSEI_VERSION)
	@ln -sf licensei-$(LICENSEI_VERSION) bin/licensei

bin/licensei-$(LICENSEI_VERSION): | $(BIN_DIR)
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v$(LICENSEI_VERSION)
	@mv bin/licensei $@

####
## ActionLint CLI
####

ACTIONLINT_BIN := $(BIN_DIR)/actionlint
ACTIONLINT_VERSION := 1.6.26

bin/actionlint: bin/actionlint-$(ACTIONLINT_VERSION)
	@ln -sf actionlint-$(ACTIONLINT_VERSION) bin/actionlint

bin/actionlint-$(ACTIONLINT_VERSION): | $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/rhysd/actionlint/main/scripts/download-actionlint.bash \
	| bash -s -- "$(ACTIONLINT_VERSION)" "$(BIN_DIR)"
	@mv bin/actionlint $@

####
##  Azure CLI
####

AZURECLI_BIN := $(BIN_DIR)/az
AZURECLI_VERSION := 2.53.0
AZURECLI_VENV := $(AZURECLI_BIN)-$(AZURECLI_VERSION)

bin/az: $(AZURECLI_VENV)/bin/az
	@ln -sf $(AZURECLI_VENV)/bin/az bin/az

$(AZURECLI_VENV)/bin/az: | $(BIN_DIR)
	@python3 -m venv $(AZURECLI_VENV)
	@$(AZURECLI_VENV)/bin/python3 -m pip install --upgrade pip
	@$(AZURECLI_VENV)/bin/pip install azure-cli==$(AZURECLI_VERSION)

####
##  Azure Bicep CLI
####

BICEP_BIN := $(BIN_DIR)/bicep
BICEP_VERSION := 0.22.6
BICEP_OSTYPE := $(OSTYPE)
BICEP_ARCH := $(ARCHTYPE)

# Set OSTYPE for macos to "osx"
ifeq ($(BICEP_OSTYPE),darwin)
	BICEP_OSTYPE = osx
endif
# Reset ARCHTYPE for amd64 to "x64"
ifeq ($(BICEP_ARCH),amd64)
	BICEP_ARCH = x64
endif

bin/bicep: bin/bicep-$(BICEP_VERSION)
	@ln -sf bicep-$(BICEP_VERSION) bin/bicep

bin/bicep-$(BICEP_VERSION): | $(BIN_DIR)
	@if [ -z "${BICEP_OSTYPE}" -o -z "${BICEP_ARCH}" ]; then printf 'ERROR: following variables must no be empty: %s %s\n' '$$BICEP_OSTYPE' '$$BICEP_ARCH'; exit 1; fi
	@curl -sSfL 'https://github.com/Azure/bicep/releases/download/v$(BICEP_VERSION)/bicep-$(BICEP_OSTYPE)-$(BICEP_ARCH)' \
	--output '$(BICEP_BIN)-$(BICEP_VERSION)'
	@chmod +x '$(BICEP_BIN)-$(BICEP_VERSION)'

####
##  CloudFormation Linter CLI
####

CFNLINT_BIN := $(BIN_DIR)/cfn-lint
CFNLINT_VERSION := 0.82.2
CFNLINT_VENV := $(CFNLINT_BIN)-$(CFNLINT_VERSION)

bin/cfn-lint: $(CFNLINT_VENV)/bin/cfn-lint
	@ln -sf $(CFNLINT_VENV)/bin/cfn-lint bin/cfn-lint

$(CFNLINT_VENV)/bin/cfn-lint: | $(BIN_DIR)
	@python3 -m venv $(CFNLINT_VENV)
	@$(CFNLINT_VENV)/bin/python3 -m pip install --upgrade pip
	@$(CFNLINT_VENV)/bin/pip install cfn-lint==$(CFNLINT_VERSION)

####
##  Golangci-lint CLI
####

GOLANGCI_BIN := $(BIN_DIR)/golangci-lint
GOLANGCI_CONFIG := $(ROOT_DIR)/.golangci.yml
GOLANGCI_VERSION := 1.54.2

bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) bin/golangci-lint

bin/golangci-lint-$(GOLANGCI_VERSION): | $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b "$(BIN_DIR)" "v$(GOLANGCI_VERSION)"
	@mv bin/golangci-lint $@
