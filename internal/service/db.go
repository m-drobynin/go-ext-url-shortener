package service

type DBProvider interface {
	TryPut(key string, value string) error
	Get(key string) (*string, error)
}
