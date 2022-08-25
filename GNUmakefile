default: testacc
PKG_NAME = sonarr

# Local provider install parameters
version = 0.1.0
registry_name = registry.terraform.io
namespace = $(PKG_NAME)
bin_name = terraform-provider-$(PKG_NAME)
build_dir = .build
TF_PLUGIN_DIR ?= ~/.local/share/terraform/plugins
install_path = $(TF_PLUGIN_DIR)/$(registry_name)/$(namespace)/$(PKG_NAME)/$(version)/$$(go env GOOS)_$$(go env GOARCH)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build plugin binary
.PHONY: build
build:
	mkdir -p $(build_dir)
	go build -tags all -o $(build_dir)/$(bin_name)

# Install the binary in the plugin directory
.PHONY: install
install: build
	mkdir -p $(install_path)
	cp $(build_dir)/$(bin_name) $(install_path)/$(bin_name)

# Generate documentation
.PHONY: doc
doc:
	go generate ./...

# Lint
.PHONY: lint
lint:
	golangci-lint run ./...

# Format
.PHONY: fmt
fmt:
	go fmt ./...
