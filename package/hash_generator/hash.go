package hash_generator

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/phamduytien1805/package/config"
	"golang.org/x/crypto/argon2"
)

// HashSalt struct used to store
// generated hash and salt used to
// generate the hash.
type HashSalt struct {
	Hash, Salt string
}

type Argon2idHash struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

func NewArgon2idHash(config *config.Config) *Argon2idHash {
	return &Argon2idHash{
		time:    config.Hash.Time,
		saltLen: config.Hash.SaltLen,
		memory:  config.Hash.Memory,
		threads: config.Hash.Threads,
		keyLen:  config.Hash.KeyLen,
	}
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	// If salt is not provided generate a salt of
	// the configured salt length.
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return &HashSalt{Hash: base64.StdEncoding.EncodeToString(hash), Salt: base64.StdEncoding.EncodeToString(salt)}, nil
}

// Compare generated hash with store hash.
func (a *Argon2idHash) Compare(hash, salt, password string) error {
	// Decode the hash and salt from base64.
	rawHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	rawSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return err
	}
	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash([]byte(password), rawSalt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(rawHash, []byte(hashSalt.Hash)) {
		return errors.New("hash doesn't match")
	}
	return nil
}
