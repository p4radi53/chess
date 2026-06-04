.PHONY: lint lint-go lint-fe fmt fmt-go fmt-fe fmt-check fmt-check-go fmt-check-fe

lint: lint-go lint-fe

lint-go:
	golangci-lint run ./internal/... ./cmd/...

lint-fe:
	cd web/chess && npm run lint && npx tsc --noEmit

fmt: fmt-go fmt-fe

fmt-go:
	gofmt -w ./internal ./cmd

fmt-fe:
	cd web/chess && npx prettier --write .

fmt-check: fmt-check-go fmt-check-fe

fmt-check-go:
	gofmt -l ./internal ./cmd | grep . && exit 1 || exit 0

fmt-check-fe:
	cd web/chess && npx prettier --check .
