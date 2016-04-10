package data

import (
	"testing"
	"fmt"
)

func TestNewRedisDatabaseManager(t *testing.T) {
	r := NewRedisDatabaseManager()

	r.Insert("foo", "bar")
	result, _ := r.Client.Get("foo").Result()
	if result != "bar" {
		fmt.Println("Key \"foo\" not found")
		t.Fail()
	}
	r.Client.Del("foo")

	r.Insert("foo", []string{"foo", "bar", "baz"})
	if r.max != 3 {
		t.Fail()
	}
	result, _ = r.Client.LPop("foo").Result()
	if result != "foo" {
		t.Fail()
	}
	result, _ = r.Client.LPop("foo").Result()
	if result != "bar" {
		t.Fail()
	}
	result, _ = r.Client.LPop("foo").Result()
	if result != "baz" {
		t.Fail()
	}
	r.Client.Del("foo")

	r.Insert("foo", map[string]string{
		"bar": "baz",
		"hello": "world",
	})
	results, _ := r.Client.HGetAll("foo").Result()
	if results[0] != "bar" || results[1] != "baz" || results[2] != "hello" || results[3] != "world" {
		t.Fail()
	}
	r.Client.Del("foo")
	var c chan string
	invalidResult, err := r.Insert("foo", c)
	if invalidResult != nil || err == nil {
		t.Fail()
	}
}

func TestRedisDatabaseManager_Update(t *testing.T) {
	r := NewRedisDatabaseManager()
	r.Insert("foo", "bar")
	r.Update("foo", "baz")
	if r.Client.Get("foo").Val() != "baz" {
		t.Fail()
	}
}

func TestRedisDatabaseManager_Search(t *testing.T) {
	r := NewRedisDatabaseManager()
	r.Client.Set("foo", "bar", 0)
	if result, _ := r.Search("foo"); result != "bar" {
		t.Fail()
	}
	r.Del("foo")

	r.Insert("foo", map[string]string{
		"bar": "baz",
		"hello": "world",
	})
	mapResult, err := r.Search("foo")
	if data, ok := mapResult.(map[string]string); ok && err == nil {
		if data["bar"] != "baz" || data["hello"] != "world" {
			t.Fail()
		}
	}
	r.Del("foo")

	r.Insert("foo", []string{"bar", "baz", "hello"})
	strListResult, _ := r.Search("foo");
	if data, ok := strListResult.([]string); ok {
		if data[0] != "bar" || data[1] != "baz" || data[2] != "hello" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
	r.Del("foo")

	if _, err := r.Search("asdoo"); err == nil {
		t.Fail()
	}
}

func TestRedisDatabaseManager_Delete(t *testing.T) {
	r := NewRedisDatabaseManager()
	r.Insert("foo", "bar")
	r.Delete("foo")
	if _, err := r.Search("foo"); err == nil {
		t.Fail()
	}
}