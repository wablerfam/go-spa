.PHONY: dev build clean wasm server

# Development mode: build wasm and run server
dev: wasm
	go run .

# Build wasm binary
wasm:
	GOARCH=wasm GOOS=js go build -o web/app.wasm .

# Production build
build: wasm
	go build -o counter-app .

# Clean build artifacts
clean:
	rm -f web/app.wasm counter-app
