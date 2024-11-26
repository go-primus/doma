package idgen

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(3)
}

func GenerateID() string {

	id := node.Generate()
	return id.String()
}
