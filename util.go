package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlDataSource = "sqlDataSource"
	redisOptions  = "redisOptions"
)

func config(file string) (*sql.DB, *redis.Client, error) {
	//Has to come from secret manager, using local config.json for now
	f, _ := os.Open(file)
	defer f.Close()

	var c map[string]interface{}
	byteValue, _ := ioutil.ReadAll(f)
	json.Unmarshal([]byte(byteValue), &c)

	db, err := db(c[sqlDataSource].(string))

	if err != nil {
		return nil, nil, err
	}

	rAddress := c[redisOptions].(map[string]interface{})["address"].(string)
	rPass := c[redisOptions].(map[string]interface{})["password"].(string)

	cache, err := cache(rAddress, rPass)
	if err != nil {
		return nil, nil, err
	}

	return db, cache, nil

}
func db(dataSource string) (*sql.DB, error) {
	sql, err := sql.Open("mysql", dataSource)

	if err != nil {
		return nil, err
	}
	if err = sql.Ping(); err != nil {
		return nil, err
	}

	return sql, nil
}

func cache(addr string, pass string) (*redis.Client, error) {

	var ctx = context.Background()
	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	if _, err := cache.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return cache, nil
}
