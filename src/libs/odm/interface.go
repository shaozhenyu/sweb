package odm

type IdMaker interface {
	Incr(key string) interface{}
}
