//go:generate mockgen -package secretsstore -source=secretsstore.go -destination secretsstore_mock.go

// Package secretsstore persists system secrets values
package secretsstore

type SecretsStorer interface {
	// Get returns a parameter value specified by name.
	// If the name does not exist, a empty string is returned
	Get(name string) (string, error)

	// Set a string value against the specified name
	Set(name, value string) error
}

type SecretsStore struct {
	storer SecretsStorer
}

func NewSecretsStore(storer SecretsStorer) *SecretsStore {
	return &SecretsStore{storer: storer}
}

func (s *SecretsStore) Get(name string) (string, error) {
	return s.storer.Get(name)
}

func (s *SecretsStore) Set(name string, value string) error {
	return s.storer.Set(name, value)
}
