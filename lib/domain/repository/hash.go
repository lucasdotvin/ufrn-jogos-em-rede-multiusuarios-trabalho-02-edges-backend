package repository

type HashRepository interface {
	Hash(plain string) (string, error)

	Compare(hashed string, plain string) bool
}
