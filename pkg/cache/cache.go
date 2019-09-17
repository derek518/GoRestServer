package cache

import (
	"errors"
	"github.com/gin-contrib/cache/persistence"
	"github.com/rs/zerolog/log"
	"time"
	"GoRestServer/pkg/config"
)

var Store *persistence.RedisStore

func Setup() {
	Store = persistence.NewRedisCache(config.Redis.Host, config.Redis.Password, config.Redis.Expiration)
}

func Set(key string, value interface{}) error {
	if *config.Redis.Enabled && Store != nil {
		err := Store.Set(key, value, time.Duration(config.Redis.Expiration*time.Second))
		log.Debug().Interface(key, value).Err(err).Msg("Set cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func Add(key string, value interface{}) error {
	if *config.Redis.Enabled && Store != nil {
		err := Store.Add(key, value, time.Duration(config.Redis.Expiration*time.Second))
		log.Debug().Interface(key, value).Err(err).Msg("Add cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func AddWithTimeout(key string, value interface{}, expires time.Duration) error {
	if *config.Redis.Enabled && Store != nil {
		if expires <= time.Duration(0) {
			expires = config.Redis.Expiration * time.Second
		}
		err := Store.Add(key, value, expires)
		log.Debug().Interface(key, value).Err(err).Msg("Add cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func Replace(key string, value interface{}) error {
	if *config.Redis.Enabled && Store != nil {
		err := Store.Replace(key, value, time.Duration(config.Redis.Expiration*time.Second))
		log.Debug().Interface(key, value).Err(err).Msg("Replace cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func Get(key string, ptrValue interface{}) error {
	if *config.Redis.Enabled && Store != nil {
		err := Store.Get(key, ptrValue)
		log.Debug().Interface(key, ptrValue).Err(err).Msg("Get cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func Delete(key string) error {
	if *config.Redis.Enabled && Store != nil {
		err := Store.Delete(key)
		log.Debug().Str("key", key).Err(err).Msg("Delete cache")
		return err
	}
	return errors.New("cache is not enabled")
}

func Increment(key string, delta uint64) (uint64, error) {
	if *config.Redis.Enabled && Store != nil {
		r, err := Store.Increment(key, delta)
		log.Debug().Int64(key, int64(r)).Err(err).Msg("Increment cache")
		return r, err
	}
	return 0, errors.New("cache is not enabled")
}

func Decrement(key string, delta uint64) (newValue uint64, err error) {
	if *config.Redis.Enabled && Store != nil {
		r, err := Store.Decrement(key, delta)
		log.Debug().Int64(key, int64(r)).Err(err).Msg("Decrement cache")
		return r, err
	}
	return 0, errors.New("cache is not enabled")
}

func Flush() error {
	if *config.Redis.Enabled && Store != nil {
		return Store.Flush()
	}
	return errors.New("cache is not enabled")
}
