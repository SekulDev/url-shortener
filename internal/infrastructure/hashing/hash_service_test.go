package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashService_GenerateHash(t *testing.T) {
	node := InitSnowflakeNode(1)
	hashService := NewHashService(node)

	hash := hashService.GenerateHash()

	assert.Regexp(t, "^[0-9A-Za-z]+$", hash, "Expected hash to be a base62 string")
}
