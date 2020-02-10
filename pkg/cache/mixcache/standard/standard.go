package standard

type MixCache interface {
	Set(key string, value interface{}, expire int) (resp interface{}, err error)
	Get(key string) (resp interface{}, err error)
}
