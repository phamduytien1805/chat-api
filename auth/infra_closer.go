package main

import (
	"github.com/phamduytien1805/auth/infrastructures/userclient"
	redis_engine "github.com/phamduytien1805/package/redis"
)

type InfraStruct struct {
}

func NewInfraCloser() *InfraStruct {
	return &InfraStruct{}
}

func (i *InfraStruct) Close() error {
	err := redis_engine.RedisConn.Close()
	if err != nil {
		return err
	}
	return userclient.UserClientConn.Close()
}
