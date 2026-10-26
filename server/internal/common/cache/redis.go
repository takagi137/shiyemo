package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shiyemo/server/internal/common/config"
)

// RedisClient Redis 客户端封装
type RedisClient struct {
	client *redis.Client
}

// 全局 Redis 客户端实例
var globalClient *redis.Client

// Init 初始化 Redis 连接
func Init(cfg config.RedisConfig) error {
	globalClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := globalClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return globalClient
}

// Close 关闭 Redis 连接
func Close() error {
	if globalClient != nil {
		return globalClient.Close()
	}
	return nil
}

// ============================================================
// 缓存工具函数
// ============================================================

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if globalClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value: %w", err)
	}
	return globalClient.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string, dest interface{}) error {
	if globalClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	data, err := globalClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	if globalClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return globalClient.Del(ctx, keys...).Err()
}

// DeletePattern 按模式删除缓存
func DeletePattern(ctx context.Context, pattern string) error {
	if globalClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	iter := globalClient.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := globalClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Exists 检查缓存是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := globalClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// ============================================================
// 缓存 Key 常量
// ============================================================

const (
	// CacheKeyPostList 帖子列表缓存 Key（带页码和类型）
	CacheKeyPostList = "shiyemo:post:list:%s:%d:%d" // type:page:pageSize
	// CacheKeyPostDetail 帖子详情缓存 Key
	CacheKeyPostDetail = "shiyemo:post:detail:%d" // postID
	// CacheKeyUserDetail 用户详情缓存 Key
	CacheKeyUserDetail = "shiyemo:user:detail:%d" // userID
)

// 缓存过期时间
const (
	CacheExpirePostList   = 5 * time.Minute  // 帖子列表 5 分钟
	CacheExpirePostDetail = 10 * time.Minute // 帖子详情 10 分钟
	CacheExpireUserDetail = 15 * time.Minute // 用户详情 15 分钟
	CacheExpireUserStats  = 5 * time.Minute  // 用户统计 5 分钟
)

// FormatPostListKey 生成帖子列表缓存 Key
func FormatPostListKey(postType string, page, pageSize int) string {
	return fmt.Sprintf(CacheKeyPostList, postType, page, pageSize)
}

// FormatPostDetailKey 生成帖子详情缓存 Key
func FormatPostDetailKey(postID int64) string {
	return fmt.Sprintf(CacheKeyPostDetail, postID)
}

// FormatUserDetailKey 生成用户详情缓存 Key
func FormatUserDetailKey(userID int64) string {
	return fmt.Sprintf(CacheKeyUserDetail, userID)
}

// ============================================================
// Singleflight 防缓存雪崩
// ============================================================

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Singleflight 轻量级请求合并器，防止缓存雪崩
type Singleflight struct {
	mu sync.Mutex
	m  map[string]*call
}

// NewSingleflight 创建 Singleflight 实例
func NewSingleflight() *Singleflight {
	return &Singleflight{m: make(map[string]*call)}
}

// Do 执行请求合并：相同 key 只触发一次 fn，其余等待结果
func (s *Singleflight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()
	if s.m == nil {
		s.m = make(map[string]*call)
	}
	if c, ok := s.m[key]; ok {
		s.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := &call{}
	c.wg.Add(1)
	s.m[key] = c
	s.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()

	return c.val, c.err
}
