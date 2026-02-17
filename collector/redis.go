package collector

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type RedisCollector struct {
	addr   string
	client *redis.Client

	upDesc               *prometheus.Desc
	connectedClientsDesc *prometheus.Desc
	usedMemoryDesc       *prometheus.Desc

	// ✨--- 新增指标描述符 --🪐---
	commandsProcessedDesc *prometheus.Desc
	keyspaceHitsDesc      *prometheus.Desc
	keyspaceMissesDesc    *prometheus.Desc
	expiredKeysDesc       *prometheus.Desc
	evictedKeysDesc       *prometheus.Desc
	// ✨---
}

func NewRedisCollector(addr string) *RedisCollector {
	return &RedisCollector{
		addr: addr,
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		upDesc: prometheus.NewDesc(
			"redis_up",
			"Whether Redis is up",
			[]string{"addr"}, // --🪐---
			nil,
		),
		connectedClientsDesc: prometheus.NewDesc(
			"redis_connected_clients",
			"Number of connected clients",
			[]string{"addr"}, // --🪐---
			nil,
		),
		usedMemoryDesc: prometheus.NewDesc(
			"redis_memory_used_bytes",
			"Memory used by Redis in bytes",
			[]string{"addr"}, // --🪐---
			nil,
		),

		// ✨--- 初始化新增指标 --🪐---
		commandsProcessedDesc: prometheus.NewDesc(
			"redis_commands_processed_total",
			"Total number of commands processed",
			[]string{"addr"}, // --🪐---
			nil,
		),
		keyspaceHitsDesc: prometheus.NewDesc(
			"redis_keyspace_hits_total",
			"Keyspace hits",
			[]string{"addr"}, // --🪐---
			nil,
		),
		keyspaceMissesDesc: prometheus.NewDesc(
			"redis_keyspace_misses_total",
			"Keyspace misses",
			[]string{"addr"}, // --🪐---
			nil,
		),
		expiredKeysDesc: prometheus.NewDesc(
			"redis_expired_keys_total",
			"Number of expired keys",
			[]string{"addr"}, // --🪐---
			nil,
		),
		evictedKeysDesc: prometheus.NewDesc(
			"redis_evicted_keys_total",
			"Number of evicted keys",
			[]string{"addr"}, // --🪐---
			nil,
		),
		// ✨---
	}
}

// Describe 告诉 Prometheus：我会吐哪些指标
func (c *RedisCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upDesc
	ch <- c.connectedClientsDesc
	ch <- c.usedMemoryDesc

	// ✨--- 新增指标 --🪐---
	ch <- c.commandsProcessedDesc
	ch <- c.keyspaceHitsDesc
	ch <- c.keyspaceMissesDesc
	ch <- c.expiredKeysDesc
	ch <- c.evictedKeysDesc
	// ✨---
}

// Collect 真正采集指标
func (c *RedisCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	// 1. 检查 Redis 是否存活
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, c.addr) // --🪐---
		return
	}

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, c.addr) // --🪐---

	// 2. 拉 INFO
	info, err := c.client.Info(ctx).Result()
	if err != nil {
		return
	}

	metrics := parseInfo(info)

	// connected_clients
	if v, ok := metrics["connected_clients"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.connectedClientsDesc, prometheus.GaugeValue, val, c.addr) // --🪐---
		}
	}

	// used_memory
	if v, ok := metrics["used_memory"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.usedMemoryDesc, prometheus.GaugeValue, val, c.addr) // --🪐---
		}
	}

	// ✨--- 新增指标采集 --🪐---
	if v, ok := metrics["total_commands_processed"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.commandsProcessedDesc, prometheus.CounterValue, val, c.addr) // --🪐---
		}
	}

	if v, ok := metrics["keyspace_hits"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.keyspaceHitsDesc, prometheus.CounterValue, val, c.addr) // --🪐---
		}
	}

	if v, ok := metrics["keyspace_misses"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.keyspaceMissesDesc, prometheus.CounterValue, val, c.addr) // --🪐---
		}
	}

	if v, ok := metrics["expired_keys"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.expiredKeysDesc, prometheus.CounterValue, val, c.addr) // --🪐---
		}
	}

	if v, ok := metrics["evicted_keys"]; ok {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(c.evictedKeysDesc, prometheus.CounterValue, val, c.addr) // --🪐---
		}
	}
	// ✨---
}

// INFO 文本解析
func parseInfo(info string) map[string]string {
	result := make(map[string]string)

	for _, line := range strings.Split(info, "\r\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}
