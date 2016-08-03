package odm

type IdMaker interface {
	Incr(key string) interface{}
}

type IBeforeInsert interface {
	BeforeInsert(db *DB) error
}
