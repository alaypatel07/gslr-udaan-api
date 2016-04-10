package data

import (
	"gopkg.in/redis.v3"
	"errors"
)

type Result interface {}

type DatabaseManager interface {
	Insert(key string, v interface{}) (Result, error)
	Update(key string, new interface{}) (Result, error)
	Delete(key string) (Result, error)
	Search(key string) (Result, error)
}

type RedisDatabaseManager struct {
	*redis.Client
	max int
}

func NewRedisDatabaseManager() *RedisDatabaseManager {
	return &RedisDatabaseManager{redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		},
	), 0}
}

func (r *RedisDatabaseManager)Insert(key string, v interface{}) (Result, error) {
	switch v.(type) {
	case string:
		data := v.(string)
		return r.Set(key, data, 0).Result()
	case []string:
		data := v.([]string)
		for _, value := range data {
			if _, err := r.RPush(key, value).Result(); err != nil {
				return nil, err
			}
		}
		if r.max < len(data) {
			r.max = len(data)
		}
		return "OK", nil
	case map[string]string:
		data := v.(map[string]string)
		for field, value := range data {
			if _, err := r.HSet(key, field, value).Result(); err != nil {
				return nil, err
			}
		}
		return "OK", nil
	default:
		return nil, errors.New("Not a valid Redis data type")
	}
	return nil, errors.New("Unexpected Error")
}

func (r *RedisDatabaseManager)Delete(key string) (Result, error) {
	return r.Del(key).Result()
}

func (r *RedisDatabaseManager)Update(key string, v interface{}) (Result, error) {
	_, err := r.Del(key).Result()
	if err != nil {
		return nil, err
	}
	return r.Insert(key, v)
}

func (r *RedisDatabaseManager) Search(key string) (Result, error) {
	if result, err := r.Get(key).Result(); err == nil {
		return result, nil
	}
	if result, err := r.HGetAll(key).Result(); err == nil && len(result) > 0 {
		temp := make(map[string]string)
		for i := 0; i < len(result); i += 2 {
			temp[result[i]] = result[i + 1]
		}
		return temp, nil
	}
	if result, err := r.LRange(key,0, int64(r.max)).Result(); err == nil && len(result) > 0 {
		return result, err
	}
	return nil, errors.New("Not Found")
}