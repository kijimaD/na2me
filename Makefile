.DEFAULT_GOAL := help

.PHONY: run
sample: ## サンプルを実行する
	go run .

.PHONY: test
test: ## テストを実行する
	go test ./... -v

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
