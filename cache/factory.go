package cache

func NewCache(useRedis bool, redisAddr string) (Cache, error) {
	if useRedis {
		return NewRedisCache(redisAddr), nil
	} else {
		return NewInMemCache(), nil
	}
}
