# HPC Wiki MCP Server

UEC HPC研究室のDokuWikiを読み取り専用で検索・取得するGo製MCPサーバーです。

## MCP tools

- `search_wiki`: Wikiを全文検索し、ページID・タイトル・抜粋・ヒット数・更新日時を返す
- `get_wiki_page`: ページIDから本文・見出し・リンクを返す

Wikiの編集、ファイルのアップロード、認証状態の変更は行いません。

## Requirements

- Go 1.26以降
- Docker
- mise

## Build and test

```bash
go mod download
go test ./...
go build -o hpc-wiki-mcp .
```

## Authentication

認証情報はファイルへ保存せず、実行環境から渡します。

```bash
export HPC_WIKI_USERNAME='<username>'
export HPC_WIKI_PASSWORD='<password>'
```

利用できる環境変数は次のとおりです。

| Variable | Default | Description |
| --- | --- | --- |
| `HPC_WIKI_BASE_URL` | `http://www.hpc.is.uec.ac.jp/wiki/hpcwiki2/doku.php` | DokuWikiエンドポイント |
| `HPC_WIKI_USERNAME` | unset | HTTP Basic認証ユーザー名 |
| `HPC_WIKI_PASSWORD` | unset | HTTP Basic認証パスワード |
| `HPC_WIKI_TIMEOUT_MS` | `15000` | Wikiへのリクエストタイムアウト |
| `PORT` | `8080` | Streamable HTTPの待受ポート |

`HPC_WIKI_USERNAME` と `HPC_WIKI_PASSWORD` は両方を設定するか、両方とも設定しないでください。

## Run over stdio

ローカルMCPクライアントから起動する場合に使用します。

```bash
./hpc-wiki-mcp serve --stdin
```

MCPクライアントの設定例:

```json
{
  "mcpServers": {
    "hpc-wiki": {
      "command": "/Users/yz/repos/hpc-wiki-mcp/hpc-wiki-mcp",
      "args": ["serve", "--stdin"],
      "env": {
        "HPC_WIKI_USERNAME": "<username>",
        "HPC_WIKI_PASSWORD": "<password>"
      }
    }
  }
}
```

## Run over Streamable HTTP

`uec-portal-mcp` と同様に、`serve` の既定モードはStreamable HTTPです。

```bash
./hpc-wiki-mcp serve
./hpc-wiki-mcp serve --http --port 18080
```

`PORT` 環境変数が設定されている場合は `--port` より優先されます。

## OAuth

HTTP モードでは、`MCP_PASSWORD` と `MCP_PUBLIC_URL` を両方設定すると [mcpauth](https://github.com/yz4230/mcpauth) による OAuth で `/mcp` を保護します。ChatGPT や Claude のコネクタから接続すると、パスワード入力ページが表示されます。どちらも未設定なら認証なしで動作し、片方だけの場合は起動に失敗します。

| 変数 | 既定値 | 説明 |
|---|---|---|
| `MCP_PASSWORD` | なし | ログインページのパスワード |
| `MCP_PUBLIC_URL` | なし | クライアントから見た公開 URL(例: `https://mcp.example.com`) |
| `MCP_STATE_FILE` | `oauth_state.json`(Docker では `/data/oauth_state.json`) | 登録済みクライアントとトークンの保存先 |
| `MCP_REDIRECT_URIS` | なし | ChatGPT と Claude に加えて許可するリダイレクト URI(カンマ区切り) |

`docker:run` は `<CONTAINER_NAME>-data` ボリュームを `/data` にマウントするので、コンテナを作り直してもトークンは保持されます。

## Run with Docker

Dockerイメージのビルドとコンテナ操作には `mise` taskを使用できます。

```bash
mise run docker:build
mise run docker:run
mise run docker:stop
mise run docker:rm
```

`docker:run` は `.env` が存在する場合にコンテナへ読み込み、既定でホストとコンテナの `8080` 番ポートを接続します。`HOST_PORT`、`CONTAINER_PORT`、`IMAGE_NAME`、`CONTAINER_NAME` で各値を変更できます。

## Project layout

```text
.
├── cmd/          # Cobra CLI and MCP transports
├── pkg/tools/    # MCP tool definitions and typed handlers
├── pkg/wiki/     # DokuWiki HTTP client and HTML parsers
├── main.go
└── go.mod
```
