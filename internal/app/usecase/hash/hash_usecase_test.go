package usecase

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"url-shortener/pkg"
)

func TestHashUsecase_GenerateHash(t *testing.T) {
	node := pkg.InitSnowflakeNode(1)
	hashService := NewHashUsecase(node)

	hash := hashService.GenerateHash()

	assert.Regexp(t, "^[0-9A-Za-z]+$", hash, "Expected hash to be a base62 string")
}
