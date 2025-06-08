# 文書翻訳CLIツール要件定義書

## 1. プロジェクト概要

### 1.1 目的

標準入出力で任意の文書を受け取り、Claude Code SDKで翻訳して元の形式を保持したまま出力するシンプルなCLIツール

### 1.2 UNIX哲学

「一つのことをうまくやる」- あらゆる形式の文書翻訳に特化、元の形式を完全保持

## 2. 機能要件

### 2.1 コア機能（必須）

- **FR-001**: 標準入力から任意の文書（Markdown、テキスト、HTML、その他）を受け取り
- **FR-002**: 指定された言語に翻訳
- **FR-003**: 翻訳結果を元の形式を保持して標準出力に出力
- **FR-004**: `claude -p`コマンドを内部で実行
- **FR-005**: 入力文書の形式（構文、構造、フォーマット）の完全保持

### 2.2 コマンドライン仕様

```bash
xlat [言語] [変換指示（オプション）]

# 基本使用例
cat README.md | xlat ja
cat article.html | xlat "スペイン語" > article_es.html
cat notes.txt | xlat en

# 変換指示付き
cat spec.md | xlat ja "技術仕様書をユーザーガイドに変換"
```

### 2.3 内部動作

1. 標準入力から文書を読み取り
2. 入力文書の形式を自動判別（Markdown、HTML、プレーンテキスト等）
3. 形式保持を含む翻訳用プロンプトを生成
4. `claude -p "プロンプト"`を実行
5. 結果を標準出力に出力

## 3. 技術要件

### 3.1 実装方針

- **言語**: シェルスクリプトまたはPython（シンプル重視）
- **依存**: Claude Code SDK（既にインストール済み前提）
- **設定**: 環境変数`ANTHROPIC_API_KEY`のみ使用

### 3.2 プロンプトテンプレート

```
以下の文書を{言語}に翻訳してください。

重要：
1. 元の文書形式（Markdown、HTML、プレーンテキスト等）を完全に保持
2. 構文、タグ、記号、構造をすべて維持
3. コードブロック、URL、技術識別子は翻訳しない
4. 文書の構造と形式を絶対に変更しない

{変換指示がある場合の追加指示}

文書：
{入力文書}
```

## 4. エラーハンドリング

### 4.1 必要最小限のエラー処理

- **E001**: 標準入力が空の場合 → エラーメッセージ表示
- **E002**: `claude`コマンドが見つからない場合 → インストール案内
- **E003**: APIキーが設定されていない場合 → 設定案内

### 4.2 エラー出力

すべてのエラーメッセージは標準エラー出力（stderr）に出力

## 5. 制約事項

### 5.1 前提条件

- Claude Code SDKがインストール済み
- `ANTHROPIC_API_KEY`環境変数が設定済み
- インターネット接続が利用可能

### 5.2 対象外機能

- 設定ファイル管理
- バッチ処理
- ログ機能
- バックアップ機能
- GUI
- 拡張機能

## 6. 実装例（シェルスクリプト）

```bash
#!/bin/bash
# xlat - 文書翻訳ツール

# 引数チェック
if [ $# -eq 0 ]; then
    echo "使用方法: xlat <言語> [変換指示]" >&2
    echo "例: cat README.md | xlat ja" >&2
    exit 1
fi

TARGET_LANG="$1"
TRANSFORM_INSTRUCTION="$2"

# 標準入力チェック
if [ -t 0 ]; then
    echo "エラー: 標準入力から文書を入力してください" >&2
    exit 1
fi

# Claude Code SDKチェック
if ! command -v claude &> /dev/null; then
    echo "エラー: claude コマンドが見つかりません" >&2
    echo "Claude Code SDKをインストールしてください" >&2
    exit 1
fi

# 標準入力を読み取り
CONTENT=$(cat)

# プロンプト生成
PROMPT="以下の文書を${TARGET_LANG}に翻訳してください。

重要：
1. 元の文書形式（Markdown、HTML、プレーンテキスト等）を完全に保持
2. 構文、タグ、記号、構造をすべて維持
3. コードブロック、URL、技術識別子は翻訳しない
4. 文書の構造と形式を絶対に変更しない"

if [ -n "$TRANSFORM_INSTRUCTION" ]; then
    PROMPT="${PROMPT}

追加指示：${TRANSFORM_INSTRUCTION}"
fi

PROMPT="${PROMPT}

文書：
${CONTENT}"

# Claude Code SDK実行
echo "$PROMPT" | claude -p
```

## 7. テスト

### 7.1 基本テスト

```bash
# Markdown翻訳
echo "# Hello World" | xlat ja

# HTML翻訳
echo "<h1>Hello World</h1>" | xlat ja

# プレーンテキスト翻訳
echo "Hello World" | xlat ja

# ファイル翻訳
cat README.md | xlat "スペイン語" > README_es.md
cat index.html | xlat ja > index_ja.html

# 変換指示付き
cat spec.md | xlat ja "開発者向けをユーザー向けに変換"
```

### 7.2 エラーテスト

```bash
# 引数なし
xlat

# パイプなし
xlat ja
```

---

**この要件定義の特徴：**

- UNIX哲学に従った単機能特化
- あらゆる文書形式に対応（Markdown、HTML、テキスト等）
- 元の形式を完全保持した翻訳
- 標準入出力の完全活用
- Claude Code SDKへの薄いラッパー
- 最小限の依存関係
- シンプルな実装
