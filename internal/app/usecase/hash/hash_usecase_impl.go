package usecase

import (
	"github.com/bwmarrin/snowflake"
	"url-shortener/pkg"
)

type hashUsecaseImpl struct {
	node *snowflake.Node
}

func NewHashUsecase(node *snowflake.Node) HashUsecase {
	return &hashUsecaseImpl{node: node}
}

func (s *hashUsecaseImpl) GenerateHash() string {
	id := s.node.Generate().Int64()
	return pkg.Base10ToBase62(id)
}
