package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

type Cacher interface {
	Set(key string, value interface{}) error
	Get(key string, dest interface{}) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cacher {
	return &cache{
		client: client,
	}
}

func (c *cache) Set(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(key, jsonValue, 0).Err()
}

func (c *cache) Get(key string, dest interface{}) error {
	result, err := c.client.Get(key).Result()

	if err == redis.Nil {
		return errors.New(fmt.Sprintf("not found by key %s", key))
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cache := NewCache(client)

	// Установка значения по ключу
	if err := cache.Set("some:key", "value"); err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	// Получение значения по ключу
	var value string
	if err := cache.Get("some:key", &value); err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}

	fmt.Println(value)

	user := &User{
		ID:   1,
		Name: "John",
		Age:  30,
	}

	// Установка значения по ключу
	if err := cache.Set(fmt.Sprintf("user:%v", user.ID), user); err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}

	// Получение значения по ключу
	var retrievedUser User
	if err := cache.Get(fmt.Sprintf("user:%v", user.ID), &retrievedUser); err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Printf("%+v\n", retrievedUser)
}
