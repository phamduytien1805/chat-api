package domain

type Hash interface {
	GenerateHash(password []byte, salt []byte) (string, error)
	Compare(hashedPassword string, password string) error
}
