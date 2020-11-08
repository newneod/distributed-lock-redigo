package main

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var conn redis.Conn

func Init(address string) {
	if address == "" {
		address = "127.0.0.1:6379"
	}

	var err error
	conn, err = redis.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Connect to redis error, %v", err)
	}
}
