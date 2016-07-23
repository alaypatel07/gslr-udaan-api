package data

import (
	"gopkg.in/redis.v3"
	"errors"
)

type Result interface {}

//DatabaseManager declares methods needed to implement for making
//any database client work with this api
type DatabaseManager interface {
	Insert(key string, v interface{}) (Result, error)
	Update(key string, new interface{}) (Result, error)
	Delete(key string) (Result, error)
	Search(key string) (Result, error)
}

//RedisDatabaseManager struct is the database client implementing
//all the methods in DatabaseManager
type RedisDatabaseManager struct {
	*redis.Client
	max int
}

//NewRedisDatabaseManager creates a new database client
func NewRedisDatabaseManager() *RedisDatabaseManager {
	//TODO: Remove the hard coding by adding necessary parameters to the function signature
	return &RedisDatabaseManager{redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		},
	), 0}
}

//Insert takes a string and a value to be inserted into the database.
//The types of value v currently supported are string, []string, map[string]string.
//If the database call to save fails or the type is not supported, it returns nil with
//appropriate error, else it returns appropriate Result, with error as nil.
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
		return nil, errors.New("This type is not supported by the api")
	}
}

//Delete takes a key and tries to delete from the database. It returns
//result and nil on success, else nil with appropriate error.
func (r *RedisDatabaseManager)Delete(key string) (Result, error) {
	return r.Del(key).Result()
}

//Update takes a key and updates it to a new value. Refer to Insert's documentation
//for more details
func (r *RedisDatabaseManager)Update(key string, v interface{}) (Result, error) {
	_, err := r.Del(key).Result()
	if err != nil {
		return nil, err
	}
	return r.Insert(key, v)
}

//Search takes a key and returns appropriate value in the as Result
//and nil. It returns nil with "Not Found" error, if the key is absent.
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
