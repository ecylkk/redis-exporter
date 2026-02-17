package main

import (
	"context" // --🌟---
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ecylkk/redis-exporter/collector"
	"github.com/redis/go-redis/v9" // --🌟---
)

var (
	// ✨--- 命令行参数 --🪐---
	listenAddr = flag.String("listen-addr", ":9121", "Address to listen on")  // HTTP metrics 端口
	redisAddr  = flag.String("redis-addr", "localhost:6379", "Redis address") // Redis 地址
	// ✨---
)

func main() {
	flag.Parse() // 解析命令行参数

	// 创建 Redis Collector，使用命令行传入的 Redis 地址 --🪐---
	redisCollector := collector.NewRedisCollector(*redisAddr)
	prometheus.MustRegister(redisCollector)

	// 暴露 /metrics
	http.Handle("/metrics", promhttp.Handler())

	// --🌟--- 新增 /health 健康检查端点 --🌟---
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		rdb := redis.NewClient(&redis.Options{Addr: *redisAddr})
		defer rdb.Close()

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable) // 503
			w.Write([]byte("Redis connection failed"))
			return
		}
		w.WriteHeader(http.StatusOK) // 200
		w.Write([]byte("OK"))
	})
	// --🌟---

	// --🌟--- 新增 / 根路径首页，方便浏览器访问 --🌟---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>Redis Exporter</title></head>
<body>
<h1>Redis Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
<p><a href="/health">Health</a></p>
</body>
</html>`))
	})
	// --🌟---

	log.Printf("Starting Redis Exporter on %s, scraping Redis at %s", *listenAddr, *redisAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
