# kc

`.env` にシークレットを書かない。Keychain に預けて、Touch ID で守る。

## 課題

- `.env` にAPIキーを平文で書いている
- `.gitignore` 頼みで、事故れば即漏洩
- チームで `.env` を Slack やメモで共有している
- AI エージェント（Claude Code, Cline 等）にシークレットへの無制限アクセスを与えたくない

## 解決

`kc` は macOS Keychain / Windows Credential Manager / Linux Secret Service をバックエンドにしたシークレット管理CLI。
`.env` には `kc://service/key` の参照だけを書き、実行時に Keychain から展開する。
シークレット取得時に Touch ID / パスワード認証が走るため、人間の承認なしにはアクセスできない。

```bash
# .env（これはgitコミットできる）
ANTHROPIC_API_KEY=kc://anthropic/api_key
AWS_SECRET_ACCESS_KEY=kc://aws/secret_access_key
PORT=3000
```

```bash
# kc 経由で起動するだけ（Touch ID が走る）
kc claude
kc npm run dev
```

## インストール

```bash
# macOS (Apple Silicon)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-darwin-arm64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc

# macOS (Intel)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-darwin-amd64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc

# Linux (amd64)
curl -fsSL https://github.com/shogomuranushi/kc/releases/latest/download/kc-linux-amd64 -o /usr/local/bin/kc && chmod +x /usr/local/bin/kc
```

## 使い方

### 1. シークレットを登録

```bash
kc set anthropic api_key sk-ant-xxxx    # 引数で渡す
kc set github token                      # プロンプト入力（非エコー）
echo "sk-xxx" | kc set stripe secret_key # パイプ
```

### 2. `.env` に参照を書く

```bash
ANTHROPIC_API_KEY=kc://anthropic/api_key
GITHUB_TOKEN=kc://github/token
PORT=3000
```

### 3. コマンドを実行

```bash
kc claude                    # .env を自動検索して展開
kc npm run dev
kc docker compose up
kc run --env-file .env.prod -- npm start  # .env ファイルを明示指定
```

### 既存の `.env` を移行する

```bash
kc migrate            # カレントディレクトリから .env を自動検索
kc migrate .env.local # ファイル指定
```

プレーンテキストの値を対話式で Keychain に移行し、`.env` を `kc://` 参照に書き換える。

### その他のコマンド

```bash
kc get github token          # 値を stdout に出力（パイプ可）
kc get github token | pbcopy # クリップボードにコピー
kc delete github token       # 削除
kc list                      # 全件一覧
kc list aws                  # サービスでフィルタ
```

## 設計

- **stdout は値のみ、メッセージはすべて stderr** → パイプで安全に使える
- **`.env` は上方向に自動検索** → サブディレクトリからでも動く
- **`kc <command>`** → `set`/`get`/`delete`/`list`/`run`/`migrate` 以外はすべて外部コマンドとして実行
- **`kc://service/key`** → service/key は `[a-zA-Z0-9_.-]` のみ許可
- **Keychain のサービス名** → `kc-cli:` プレフィックスで名前空間を分離

## 脅威モデル

| 防げること | 防げないこと |
|---|---|
| `.env` の静的スキャンによる漏洩 | プロセスに渡った後の環境変数の読み取り |
| 人間の承認なしにシークレットを取り出すこと | 実行中プロセス内での悪意あるコードによる読み取り |
| `.env` をリポジトリにコミットしての漏洩 | |

## ビルド

```bash
go build -o kc .
```
