package main

import (
	"fmt"
	"url-shortener/internal/infrastructure/hashing"
)

func main() {
	//@TODO make it everything in config bootstrap
	node := hashing.InitSnowflakeNode(1)
	hashService := hashing.NewHashService(node)
	id := hashService.GenerateHash()
	fmt.Printf("Base62  ID: %s\n", id)
}
