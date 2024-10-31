package usecase

import (
	"github.com/bwmarrin/snowflake"
	"url-shortener/pkg"
)

type HashUsecase struct {
	node *snowflake.Node
}

func NewHashUsecase(node *snowflake.Node) HashUsecase {
	return HashUsecase{node: node}
}

func (s *HashUsecase) GenerateHash() string {
	id := s.node.Generate().Int64()
	return pkg.Base10ToBase62(id)
}
