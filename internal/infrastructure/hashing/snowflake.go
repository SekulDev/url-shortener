package hashing

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

func InitSnowflakeNode(nodeID int64) *snowflake.Node {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed to initialize Snowflake node: %v", err)
	}
	return node
}
