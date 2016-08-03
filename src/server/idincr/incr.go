package idincr

import (
	"sync"
	"time"

	"libs/odm"

	//"github.com/qiniu/log"
	"gopkg.in/mgo.v2"
)

type M map[string]interface{}

type counter struct {
	Id     string `json:"id" bson:"_id"`
	Offset int64  `json:"offset" bson:"offset"`
}

type IntIDMaker struct {
	ch   map[string]chan int64
	coll *mgo.Collection
	sync.Mutex
}

func NewIntIDMaker(coll *mgo.Collection) *IntIDMaker {
	return &IntIDMaker{
		coll: coll,
		ch:   map[string]chan int64{},
	}
}

func (this *IntIDMaker) incrLoop(key string, ch_ chan int64) {
	for {
		counter_ := counter{Offset: 1}
		selector := M{"_id": key}
		chg := mgo.Change{
			Update:    M{"$inc": M{"offset": 1}},
			Upsert:    true,
			ReturnNew: true,
		}

		coll := odm.NewCollection(this.coll)
		_, err := coll.Find(selector).Apply(chg, &counter_)
		coll.Close()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		ch_ <- counter_.Offset
	}

}

func (this *IntIDMaker) Incr(key string) interface{} {
	ch, ok := this.ch[key]
	if !ok {
		ch = make(chan int64)
		this.Lock()
		this.ch[key] = ch
		this.Unlock()
		go this.incrLoop(key, ch)
	}
	return <-ch
}
