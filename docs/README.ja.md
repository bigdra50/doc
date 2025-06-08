# doc

[![Go](https://github.com/bigdra50/doc/actions/workflows/go.yml/badge.svg)](https://github.com/bigdra50/doc/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bigdra50/doc)](https://goreportcard.com/report/github.com/bigdra50/doc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

複数のLLMプロバイダー（Claude Code、OpenAI、Anthropic）を使用して、元のフォーマットを保持しながらドキュメントを翻訳するためのシンプルなコマンドラインツールです。

## 機能

- **フォーマット保持**: Markdown、HTML、プレーンテキストなどの構造を完全に維持します。
- **複数プロバイダーサポート**: Claude Code CLI（デフォルト）、OpenAI API、Anthropic API
- **35以上の言語サポート**: 日本語、英語、中国語、ロシア語など、幅広い言語をサポートしています。
- **インテリジェントなレスポンス処理**: JSON構造のレスポンス、エラーハンドリング
- **進捗表示**: アニメーションスピナーと経過時間の表示
- **シェル統合**: UNIXパイプラインの完全サポート

## インストール

```bash
# GitHubからインストール
go install github.com/bigdra50/doc@latest

# またはソースからビルド
git clone https://github.com/bigdra50/doc.git
cd doc
go build -o doc .
```

## 使用法

### 基本的な翻訳

```bash
# 標準入力から日本語に翻訳
cat document.md | doc ja

# ファイルからロシア語に翻訳（詳細なログ付き）
cat spec.html | doc -v ru

# サポートされている言語コードのリストを表示
doc --list

# カスタム指示による翻訳
cat technical_doc.md | doc ja "技術仕様をユーザーガイドに変換"
```

### プロバイダー設定

#### 1. Claude Code CLI（デフォルト）

```bash
# インストール（npmが必要）
npm install -g @anthropic-ai/claude-code

# 使用例
cat document.md | doc ja
```

#### 2. OpenAI API

```bash
# APIキーを設定
export OPENAI_API_KEY=sk-your-openai-api-key
export LLM_PROVIDER=openai

# 使用例
cat document.md | doc ja
```

#### 3. Anthropic API

```bash
# APIキーを設定
export ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key
export LLM_PROVIDER=anthropic

# 使用例（近日公開予定）
cat document.md | doc ja
```

## 設定

### 設定ファイル

ツールはTOMLファイルを介した永続的な設定をサポートしています：

```bash
# 設定ファイルを初期化
doc --init-config

# 現在の設定を表示
doc --config

# 設定値を設定
doc --set provider=openai
doc --set openai_api_key=sk-your-key
doc --set openai_model=gpt-4o
```

設定ファイルの場所はXDGベースディレクトリ仕様に従います：

- `$XDG_CONFIG_HOME/bigdra50/doc/config.toml`
- `~/.config/bigdra50/doc/config.toml`（フォールバック）

### 環境変数

環境変数は設定ファイルの設定を上書きします：

#### 環境設定ファイル

```bash
# .envファイルを作成
echo "LLM_PROVIDER=openai" > .env
echo "OPENAI_API_KEY=sk-your-api-key" >> .env
echo "OPENAI_MODEL=gpt-4o-mini" >> .env
```

### モデル選択

```bash
# 利用可能なモデルのリストを表示
doc --list-models

# プロバイダーごとのモデルを表示
doc --list-models openai
doc --list-models anthropic
```

## サポートされている言語

| コード | 言語         | コード | 言語       | コード | 言語         |
| ------ | ------------ | ------ | ---------- | ------ | ------------ |
| ja     | 日本語       | en     | 英語       | ko     | 韓国語       |
| zh     | 中国語       | ru     | ロシア語   | es     | スペイン語   |
| fr     | フランス語   | de     | ドイツ語   | it     | イタリア語   |
| pt     | ポルトガル語 | ar     | アラビア語 | hi     | ヒンディー語 |

完全なリストは `doc --list` で確認できます。

## 環境変数

| 変数名              | 説明              | デフォルト値                |
| ------------------- | ----------------- | --------------------------- |
| `LLM_PROVIDER`      | プロバイダー選択  | `claude-code`               |
| `OPENAI_API_KEY`    | OpenAI APIキー    | -                           |
| `ANTHROPIC_API_KEY` | Anthropic APIキー | -                           |
| `OPENAI_MODEL`      | OpenAIモデル      | `gpt-4o-mini`               |
| `ANTHROPIC_MODEL`   | Anthropicモデル   | `claude-3-5-haiku-20241022` |
| `CLAUDE_MODEL`      | Claude Codeモデル | `sonnet`                    |

## エラーコード

| コード | 説明                                       |
| ------ | ------------------------------------------ |
| 0      | 成功                                       |
| 1      | システムエラー（入力なし、設定の問題など） |
| 2      | 同じ言語（すでにターゲット言語）           |
| 3      | 翻訳不可能（コード、データなど）           |
| 4      | フォーマットエラー                         |
| 5      | コンテンツエラー                           |

## 実行例

```bash
# 基本的な翻訳
echo "Hello World" | doc ja
# → こんにちは世界

# 同じ言語の検出
echo "こんにちは" | doc ja
# → 終了コード 2

# Markdownフォーマットの保持
echo "# Title\n- List item" | doc ja
# → # タイトル\n- リスト項目

# 進捗表示
cat large_document.md | doc -v ja
# [INFO] ドキュメントを読み込み中...
# ⠋ Claude Code CLIで翻訳中... (2.3s)
# ✓ 翻訳完了 (2.3s)
```

## 開発とテスト

```bash
# テストを実行
go test

# ビルド
go build -o doc .

# クリーンアップ
rm -f doc
```

## アーキテクチャ

### コアコンポーネント

- **main.go**: CLI処理、アプリケーションロジック
- **provider.go**: LLMプロバイダーインターフェース
- **claude_provider.go**: Claude Code CLI実装
- **openai_provider.go**: OpenAI API実装
- **models.go**: モデルカタログ、コスト計算

### 設計原則

1. **インターフェース中心**: 統一されたLLMProvider抽象
2. **環境駆動**: `.env`ファイル、環境変数のサポート
3. **エラーハンドリング**: 詳細な終了コードシステム
4. **フォーマット保持**: 元のドキュメント構造の完全な維持

## ライセンス

MITライセンス

