BINARY_NAME = xfce4-night-mode
INSTALL_DIR = $(HOME)/.local/bin
SERVICE_DIR = $(HOME)/.config/systemd/user
VERSION     = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build install uninstall enable disable clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/xfce4-night-mode/

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	@mkdir -p $(SERVICE_DIR)
	sed 's|/usr/bin|$(INSTALL_DIR)|g' $(BINARY_NAME).service > $(SERVICE_DIR)/$(BINARY_NAME).service
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Run 'make enable' to start the service"

uninstall: disable
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	rm -f $(SERVICE_DIR)/$(BINARY_NAME).service
	@echo "Uninstalled"

enable:
	systemctl --user daemon-reload
	systemctl --user enable --now $(BINARY_NAME).service
	@echo "Service enabled and started"

disable:
	-systemctl --user disable --now $(BINARY_NAME).service
	@echo "Service disabled"

clean:
	rm -f $(BINARY_NAME)
