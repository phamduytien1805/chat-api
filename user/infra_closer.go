package main

import (
	"github.com/phamduytien1805/user/infras/db"
)

type InfraStruct struct {
}

func NewInfraCloser() *InfraStruct {
	return &InfraStruct{}
}

func (i *InfraStruct) Close() error {
	db.PGConn.Close()
	return nil
}
