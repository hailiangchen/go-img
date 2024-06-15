package cache

type ICache interface {
	Set(key string, value interface{}) error
	Get(key string) ([]byte, error)
	Del(key string) error
}
