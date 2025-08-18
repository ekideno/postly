package utils

import (
	"log"

	"github.com/bwmarrin/snowflake"
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

func GenerateSnowflakeID() string {
	return node.Generate().String()
}
