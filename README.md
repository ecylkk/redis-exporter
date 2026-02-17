# Prometheus 向け Redis Exporter

Goで実装された **本番対応の Redis Exporter** です。  
Prometheus との統合に最適化され、Redis の主要メトリクスを収集して HTTP で公開します。  
軽量で高性能、複数インスタンス監視やヘルスチェック対応済みです。

---

## 特徴

- **本番対応の Go 実装**、依存関係最小  
- Redis の主要メトリクスを収集:
  - `redis_up` – Redis の稼働状況  
  - `redis_connected_clients` – 接続中クライアント数  
  - `redis_memory_used_bytes` – メモリ使用量（バイト単位）  
  - `redis_commands_processed_total` – 処理されたコマンド総数  
  - `redis_keyspace_hits_total` / `redis_keyspace_misses_total`  
  - `redis_expired_keys_total` / `redis_evicted_keys_total`  
- 複数 Redis インスタンス対応の **addrラベル**  
- ロードバランサー用の `/health` エンドポイント  
- コマンドラインで Redis アドレスや HTTP ポートを指定可能  
- コンテナや Kubernetes で簡単にデプロイ可能  

---

## クイックスタート

```bash
# デフォルト設定で実行
go run main.go

# RedisアドレスとHTTPポートを指定
go run main.go --redis-addr=127.0.0.1:6379 --listen-addr=:9121
