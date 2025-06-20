package utils

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var node *snowflake.Node

func InitSnowflake(machineID int64) {
	var err error
	snowflake.Epoch = int64(1750366800)
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		log.Fatalf("failed to initialize snowflake: %v", err)
	}
}

func GenerateID() string {
	return node.Generate().String()
}
