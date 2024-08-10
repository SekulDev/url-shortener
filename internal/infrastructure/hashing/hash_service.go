package hashing

import (
	"github.com/bwmarrin/snowflake"
)

type HashService struct {
	node *snowflake.Node
}

func NewHashService(node *snowflake.Node) *HashService {
	return &HashService{node: node}
}

func (s *HashService) GenerateHash() string {
	id := s.node.Generate().Int64()
	return Base10ToBase62(id)
}
