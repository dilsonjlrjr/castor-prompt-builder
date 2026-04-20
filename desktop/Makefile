APP     := castor-builder
OUTDIR  := ../dist
WAILS   := $(shell go env GOPATH)/bin/wails

.PHONY: all macos windows linux clean

all: macos windows linux

# ── macOS ────────────────────────────────────────────────────────────────────
# Gera universal binary (arm64 + amd64) → .app + .dmg se tiver create-dmg
macos:
	@echo "▶ Building macOS universal..."
	@mkdir -p $(OUTDIR)
	$(WAILS) build -clean -platform darwin/universal -o $(APP)
	@cp -r build/bin/$(APP).app $(OUTDIR)/$(APP)-macos.app
	@echo "✓ macOS → $(OUTDIR)/$(APP)-macos.app"
	@if command -v create-dmg >/dev/null 2>&1; then \
		create-dmg \
			--volname "CASTOR Builder" \
			--window-size 540 380 \
			--icon-size 128 \
			--icon "$(APP).app" 160 180 \
			--hide-extension "$(APP).app" \
			--app-drop-link 380 180 \
			"$(OUTDIR)/$(APP)-macos.dmg" \
			"$(OUTDIR)/$(APP)-macos.app"; \
		echo "✓ DMG → $(OUTDIR)/$(APP)-macos.dmg"; \
	else \
		echo "ℹ  Instale create-dmg para gerar .dmg: brew install create-dmg"; \
	fi

# ── Windows ──────────────────────────────────────────────────────────────────
# Requer: brew install mingw-w64
windows:
	@echo "▶ Building Windows amd64..."
	@mkdir -p $(OUTDIR)
	@if ! command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then \
		echo "✗ mingw-w64 não encontrado. Execute:  brew install mingw-w64"; \
		exit 1; \
	fi
	$(WAILS) build -clean -platform windows/amd64 -o $(APP)
	@cp build/bin/$(APP).exe $(OUTDIR)/$(APP)-windows-amd64.exe
	@echo "✓ Windows → $(OUTDIR)/$(APP)-windows-amd64.exe"

# ── Linux ────────────────────────────────────────────────────────────────────
# Usa Docker com a imagem oficial do Wails (inclui WebKit2GTK)
linux:
	@echo "▶ Building Linux amd64 (Docker)..."
	@mkdir -p $(OUTDIR)
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "✗ Docker não encontrado."; exit 1; \
	fi
	docker run --rm \
		-v "$$(pwd)/..:/app" \
		-w /app/desktop \
		ghcr.io/wailsapp/wails-linux-builder:latest \
		wails build -clean -platform linux/amd64 -o $(APP)
	@cp build/bin/$(APP) $(OUTDIR)/$(APP)-linux-amd64
	@echo "✓ Linux → $(OUTDIR)/$(APP)-linux-amd64"

# ── Limpeza ───────────────────────────────────────────────────────────────────
clean:
	rm -rf build/bin $(OUTDIR)
