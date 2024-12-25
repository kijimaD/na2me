.DEFAULT_GOAL := help

.PHONY: run
sample: ## ゲームを実行する
	go run .

.PHONY: test
test: ## テストを実行する
	go test ./... -v

.PHONY: genbg
genbg: ## 元ファイルから背景ファイルを生成する。サイズ変更、フィルタ適用を行う
	go run . normalizeBg
	docker build . --target filter -t filter
	docker run --rm -v $(PWD):/work -w /work filter /bin/sh -c "python ./scripts/filter/filter.py"

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
