package service

type DbProvider interface {
	TryPut(key string, value string) error
	Get(key string) (*string, error)
}
