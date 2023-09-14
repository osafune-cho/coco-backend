# coco-backend
[![CI](https://github.com/osafune-cho/coco-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/osafune-cho/coco-backend/actions/workflows/ci.yml)
[![deploy](https://github.com/osafune-cho/coco-backend/actions/workflows/deploy.yml/badge.svg)](https://github.com/osafune-cho/coco-backend/actions/workflows/deploy.yml)
[![Generate Documents](https://github.com/osafune-cho/coco-backend/actions/workflows/generate-documents.yml/badge.svg)](https://github.com/osafune-cho/coco-backend/actions/workflows/generate-documents.yml)

GC59201 情報メディア創成特別講義Bで開発しているサービスのバックエンドです。

- [API server documentation](https://osafune-cho.github.io/coco-backend/)

## 開発
開発には以下のソフトウェアが必要です。

- Nix
- Niv
- Docker Compose

以下のコマンドで開発に必要なソフトウェアが入った環境に入ることができます。

```
nix-shell
```

## 実行
以下のコマンドを実行することで`result/bin/`以下にバックエンドの単一バイナリを得られます。

```
nix-build
```

実行時には環境変数の設定と`poppler-utils`に依存していることに注意してください。
必要な環境変数については`.envrc.example`を参考にしてください。

また、単に実行するだけであれば以下のコマンドを実行することでバックエンドのDockerイメージを入手できます。

```
docker pull ghcr.io/osafune-cho/coco:latest
```

設定については[docker-compose.yml](https://github.com/osafune-cho/coco-infrastructure/blob/main/docker/docker-compose.yml)を参考にしてください。

## ライセンス
- MIT
- Apache-2.0
