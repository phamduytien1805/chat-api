package hash

import (
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
	"github.com/phamduytien1805/user/domain"
)

type HashGen struct {
	engine *hash_generator.Argon2idHash
}

func NewHash(config *config.HashConfig) domain.Hash {
	hashGen := hash_generator.NewArgon2idHash(config)

	return &HashGen{
		engine: hashGen,
	}
}

func (h *HashGen) GenerateHash(password []byte, salt []byte) (string, error) {
	return h.engine.GenerateHash(password, salt)
}

func (h *HashGen) Compare(hashedPassword, password string) error {
	return h.engine.Compare(hashedPassword, password)
}
