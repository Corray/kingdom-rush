# kingdom-rush Makefile
#
# build       — ebiten desktop binary (default)
# build-term  — V1.7 terminal binary (tcell, build tag `term`)
# wasm        — WASM build + 拷贝 wasm_exec.js 到 web/
# serve       — wasm + 本地 HTTP 服务 (port 8080)
# test        — go test ./...
# clean       — 删 binary + WASM artifacts

GOROOT := $(shell go env GOROOT)
WASM_EXEC_SRC := $(GOROOT)/lib/wasm/wasm_exec.js

.PHONY: build build-term wasm serve test vet clean all

all: build build-term

build:
	go build .

build-term:
	go build -tags term -o kingdom-rush-term .

wasm: web/wasm_exec.js web/kingdom-rush.wasm

web/wasm_exec.js: $(WASM_EXEC_SRC)
	@mkdir -p web
	@if [ ! -f "$(WASM_EXEC_SRC)" ]; then \
		echo "ERROR: wasm_exec.js not found at $(WASM_EXEC_SRC)"; \
		echo "Go version >= 1.22 expected (lib/wasm/wasm_exec.js)."; \
		echo "Older Go: try \$$GOROOT/misc/wasm/wasm_exec.js"; \
		exit 1; \
	fi
	cp $(WASM_EXEC_SRC) web/wasm_exec.js

web/kingdom-rush.wasm: $(wildcard *.go) assets/levels.yaml
	@mkdir -p web
	GOOS=js GOARCH=wasm go build -o web/kingdom-rush.wasm .
	@echo "WASM built: web/kingdom-rush.wasm"
	@echo "Serve with: make serve"

serve: wasm
	@echo "Serving at http://localhost:8080 (Ctrl+C to stop)"
	@cd web && python3 -m http.server 8080

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f kingdom-rush kingdom-rush-term
	rm -f web/kingdom-rush.wasm web/wasm_exec.js
