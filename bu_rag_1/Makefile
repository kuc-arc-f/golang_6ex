.PHONY: all plugin main run clean

# デフォルト: プラグインとメインプログラムの両方をビルド
all: plugin main

# Go プラグイン (.so) のビルド
# C++ ソースコードは cgo によって自動的にコンパイルされる
plugin:
	@echo "🔨 プラグインをビルド中..."
	cd myplugin && go build -buildmode=plugin -o myplugin.so .
	@echo "✅ myplugin/myplugin.so を生成しました"

# メインプログラムのビルド
main:
	@echo "🔨 メインプログラムをビルド中..."
	go build -o app main.go
	@echo "✅ app を生成しました"

# ビルドして実行
run: all
	@echo ""
	@echo "🚀 実行中..."
	@echo ""
	./app

# クリーンアップ
clean:
	@echo "🧹 クリーンアップ中..."
	rm -f app myplugin/myplugin.so
	@echo "✅ 完了"
