APP    := castor-builder
OUTDIR := ../dist
WAILS  := $(shell command -v wails 2>/dev/null || echo "$(shell go env GOPATH)/bin/wails")

.PHONY: all macos windows linux clean generate-icns

all: macos windows linux

# ── Gera .icns com múltiplos tamanhos a partir de build/appicon.png ──────────
# Requer: sips + iconutil (nativos no macOS)
generate-icns:
	@echo "▶ Gerando castor-builder.icns a partir de appicon.png..."
	@ICONSET=build/darwin/appicon.iconset; \
	rm -rf $$ICONSET && mkdir -p $$ICONSET; \
	sips -z 16   16   build/appicon.png --out $$ICONSET/icon_16x16.png       >/dev/null; \
	sips -z 32   32   build/appicon.png --out $$ICONSET/icon_16x16@2x.png    >/dev/null; \
	sips -z 32   32   build/appicon.png --out $$ICONSET/icon_32x32.png       >/dev/null; \
	sips -z 64   64   build/appicon.png --out $$ICONSET/icon_32x32@2x.png    >/dev/null; \
	sips -z 128  128  build/appicon.png --out $$ICONSET/icon_128x128.png     >/dev/null; \
	sips -z 256  256  build/appicon.png --out $$ICONSET/icon_128x128@2x.png  >/dev/null; \
	sips -z 256  256  build/appicon.png --out $$ICONSET/icon_256x256.png     >/dev/null; \
	sips -z 512  512  build/appicon.png --out $$ICONSET/icon_256x256@2x.png  >/dev/null; \
	sips -z 512  512  build/appicon.png --out $$ICONSET/icon_512x512.png     >/dev/null; \
	sips -z 1024 1024 build/appicon.png --out $$ICONSET/icon_512x512@2x.png  >/dev/null; \
	iconutil -c icns $$ICONSET -o build/darwin/castor-builder.icns; \
	rm -rf $$ICONSET
	@echo "✓ icns → build/darwin/castor-builder.icns"

# ── macOS ────────────────────────────────────────────────────────────────────
# Gera universal binary (arm64 + amd64) → .app + .dmg se tiver create-dmg
macos: generate-icns
	@echo "▶ Building macOS universal..."
	@mkdir -p $(OUTDIR)
	$(WAILS) build -clean -platform darwin/universal -o $(APP)
	@rm -rf "$(OUTDIR)/$(APP)-macos.app"
	@cp -r "build/bin/CASTOR Builder.app" "$(OUTDIR)/$(APP)-macos.app"
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
	@cp "build/bin/$(APP)" $(OUTDIR)/$(APP)-windows-amd64.exe
	@echo "✓ Windows → $(OUTDIR)/$(APP)-windows-amd64.exe"

# ── Linux ────────────────────────────────────────────────────────────────────
# Usa Docker com a imagem oficial do Wails (inclui WebKit2GTK)
linux:
	@echo "▶ Building Linux amd64 (Docker)..."
	@mkdir -p $(OUTDIR)
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "✗ Docker não encontrado."; exit 1; \
	fi
	docker run --rm --platform linux/amd64 \
		-v "$$(pwd)/..:/app" \
		-w /app/desktop \
		castor-linux-builder \
		build -platform linux/amd64 -s -o $(APP)
	@cp build/bin/$(APP) $(OUTDIR)/$(APP)-linux-amd64
	@echo "✓ Linux → $(OUTDIR)/$(APP)-linux-amd64"

# ── Limpeza ───────────────────────────────────────────────────────────────────
clean:
	rm -rf build/bin $(OUTDIR)
