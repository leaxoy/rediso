package rediso

// RedisPool is pool of redis client.
type RedisPool struct {
	size  int
	conns []*Client
}

// NewRedisPool create new redis conn pool.
func NewRedisPool(config *Config, size int) *RedisPool {
	pool := new(RedisPool)
	pool.size = size
	for i := size; i > 0; i-- {
		redis := NewRedis(config)
		pool.conns[i] = redis
	}
	return pool
}
